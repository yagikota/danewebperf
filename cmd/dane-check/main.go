package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/yagikota/danewebperf/cmd/dane-check/model"
)

const (
	awsProfile = "default"
	awsRegion  = "ap-northeast-1"
	s3Bucket   = "pageloadtime-results"

	outPutFilePath = "../../analysis/dane-validation-all-success.csv"
)

var (
	logger         *slog.Logger
	measurementIDs = []string{
		"tokyo-v2-01",
		"tokyo-v2-02",
		"tokyo-v2-03",
		"tokyo-v2-04",
		"tokyo-v2-05",
		// "tokyo-v2-06",
		"tokyo-v2-07",
		"tokyo-v2-08",
		"tokyo-v2-09",
		"tokyo-v2-10",
		"frankfurt-v2-01",
		"frankfurt-v2-02",
		"frankfurt-v2-03",
		"frankfurt-v2-04",
		"frankfurt-v2-05",
		// "frankfurt-v2-06",
		"frankfurt-v2-07",
		"frankfurt-v2-08",
		"frankfurt-v2-09",
		"frankfurt-v2-10",
	}
)

func ListAllObjects(svc *s3.S3, bucket, prefix string) ([]*s3.Object, error) {
	var objects []*s3.Object
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Prefix: aws.String(prefix),
		Bucket: aws.String(bucket),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		objects = append(objects, p.Contents...)
		return true
	})
	if err != nil {
		return nil, err
	}
	return objects, nil
}

type CustomRetryer struct {
	client.DefaultRetryer
}

type temporary interface {
	Temporary() bool
}

func (r CustomRetryer) ShouldRetry(req *request.Request) bool {
	if origErr := req.Error; origErr != nil {
		switch origErr.(type) {
		case temporary:
			if strings.Contains(origErr.Error(), "read: connection reset") {
				// デフォルトのSDKではリトライしないが、リトライ可にする
				return true
			}
		}
	}
	return r.DefaultRetryer.ShouldRetry(req)
}

// tokyo-01/example.com/letsdane-example.com-with-cache-with-dane.csv
func targetLetsDANECSV(s3ObjectKey string) (bool, string) {
	splitKey := strings.Split(s3ObjectKey, "/")
	if len(splitKey) != 3 {
		return false, ""
	}

	if !strings.HasPrefix(splitKey[2], "letsdane-") || !strings.HasSuffix(splitKey[2], ".csv") {
		return false, ""
	}

	return true, splitKey[2]
}

func targetLetsHARCSV(s3ObjectKey string) (bool, string) {
	splitKey := strings.Split(s3ObjectKey, "/")
	if len(splitKey) != 3 {
		return false, ""
	}

	if !strings.HasSuffix(splitKey[2], ".csv") {
		return false, ""
	}

	return true, splitKey[2]
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

type Result []ResultRecord

type ResultRecord struct {
	MeasurementID string
	MeasurementInfo
	DANESuccessCount int
	TotalCount       int
	DANEAllSuccess   bool
}

type MeasurementInfo struct {
	Domain string
	Cache  string
	Dane   string
}

func newResultRecord() *ResultRecord {
	return &ResultRecord{
		MeasurementID: "",
		MeasurementInfo: MeasurementInfo{
			Domain: "",
			Cache:  "",
			Dane:   "",
		},
		DANESuccessCount: 0,
		TotalCount:       0,
		DANEAllSuccess:   false,
	}
}

func (r *ResultRecord) setMeasurementID(measurementID string) {
	r.MeasurementID = measurementID
}

func (r *ResultRecord) setDomain(domain string) {
	r.Domain = domain
}

func (r *ResultRecord) setCache(cache string) {
	r.Cache = cache
}

func (r *ResultRecord) setDane(dane string) {
	r.Dane = dane
}

func (r *ResultRecord) setMeasurementInfoFromFile(file string) {
	if strings.Contains(file, "with-cache-with-dane") {
		r.setDomain(strings.Split(file, "-with-cache-with-dane")[0])
		r.setCache("true")
		r.setDane("true")
		return
	}
	if strings.Contains(file, "with-cache-without-dane") {
		r.setDomain(strings.Split(file, "-with-cache-without-dane")[0])
		r.setCache("true")
		r.setDane("false")
		return
	}
	if strings.Contains(file, "without-cache-with-dane") {
		r.setDomain(strings.Split(file, "-without-cache-with-dane")[0])
		r.setCache("false")
		r.setDane("true")
		return
	}
	if strings.Contains(file, "without-cache-without-dane") {
		r.setDomain(strings.Split(file, "-without-cache-without-dane")[0])
		r.setCache("false")
		r.setDane("false")
		return
	}
}

func exportResultAsCSV(result Result, filPath string) error {
	file, err := os.OpenFile(filPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ファイルが空の場合のみヘッダを書き込む
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		if err := writer.Write([]string{"measurementID", "domain", "cache", "dane", "dane-success-count", "total", "dane-all-success"}); err != nil {
			return err
		}
	}

	for _, r := range result {
		csvRecord := []string{
			r.MeasurementID,
			r.Domain,
			r.Cache,
			r.Dane,
			strconv.Itoa(r.DANESuccessCount),
			strconv.Itoa(r.TotalCount),
			strconv.FormatBool(r.DANEAllSuccess),
		}
		if err := writer.Write(csvRecord); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	start := time.Now()
	log.Println("start time: ", start.Format("2006-01-02-15-04-05"))

	measurementID := flag.String("measurementID", "", "measurementID")
	outPutFilePath := flag.String("outPutFilePath", outPutFilePath, "outPutFilePath")
	flag.Parse()

	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create an AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           awsProfile,
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess, &aws.Config{
		Region:  aws.String(awsRegion),
		Retryer: CustomRetryer{},
	})

	var result Result

	prefix := *measurementID + "/"
	logger.Info(fmt.Sprintf("listing items in bucket %s with prefix %s", s3Bucket, prefix))
	objs, err := ListAllObjects(svc, s3Bucket, prefix)
	if err != nil {
		logger.Error("unable to list items in bucket %s, %v", s3Bucket, err)
		return
	}

	for idx, obj := range objs {
		logger.Info(fmt.Sprintf("now processing %s (%d/%d) %f%%", *obj.Key, idx+1, len(objs), float64(idx+1)/float64(len(objs))*100))
		ok, targetCSV := targetLetsDANECSV(*obj.Key)
		if !ok {
			continue
		}

		logger.Info(fmt.Sprintf("targetCSV: %s", targetCSV))

		gotObj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    obj.Key,
		})
		if err != nil {
			logger.Error(fmt.Sprintf("unable to get object %q from bucket %q, %v", *obj.Key, s3Bucket, err))
			continue
		}

		reader := csv.NewReader(gotObj.Body)
		records, err := reader.ReadAll()
		if err != nil {
			logger.Error(fmt.Sprintf("unable to read csv file %q, %v", *obj.Key, err))
			continue
		}

		convertedRecords := model.FilterDANEValidated(model.ConvertLetsDANERecords(records))

		if len(convertedRecords) == 0 {
			logger.Info(fmt.Sprintf("no records in %s", *obj.Key))
			continue
		}

		// dictionary to store the dane validated hosts
		dict := make(map[string][]string)
		for _, record := range convertedRecords {
			filename := getFileNameWithoutExt(targetCSV)
			key := strings.Replace(filename, "letsdane", *measurementID, 1)
			dict[key] = append(dict[key], record.Host)
		}

		s3Key := strings.Replace(*obj.Key, "letsdane-", "", 1)
		gotObj, err = svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			logger.Error(fmt.Sprintf("unable to get object %q from bucket %q, %v", s3Key, s3Bucket, err))
			continue
		}

		reader = csv.NewReader(gotObj.Body)
		records, err = reader.ReadAll()
		if err != nil {
			logger.Error(fmt.Sprintf("unable to read csv file %q, %v", s3Key, err))
			continue
		}

		convertedHarRecords := model.ConvertHarStruct(records)

		// check if all domains in convertedRecords are in the dictionary
		DANESuccessCount := 0
		for _, record := range convertedHarRecords {
			_, filename := targetLetsHARCSV(s3Key)
			key := *measurementID + "-" + getFileNameWithoutExt(filename)
			if ok := slices.Contains(dict[key], record.Domain); ok {
				DANESuccessCount++
			}
		}

		resultRecord := newResultRecord()
		resultRecord.setMeasurementID(*measurementID)
		resultRecord.setMeasurementInfoFromFile(strings.Replace(targetCSV, "letsdane-", "", 1))
		resultRecord.DANESuccessCount = DANESuccessCount
		resultRecord.TotalCount = len(convertedHarRecords)
		resultRecord.DANEAllSuccess = DANESuccessCount == len(convertedHarRecords)

		result = append(result, *resultRecord)
	}

	if err := exportResultAsCSV(result, *outPutFilePath); err != nil {
		logger.Error(fmt.Sprintf("unable to export result as csv, %v", err))
	}

	logger.Info(fmt.Sprintf("elapsed time: %f min", time.Since(start).Minutes()))
	logger.Info("done")
}
