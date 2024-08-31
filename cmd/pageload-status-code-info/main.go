package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	awsProfile = "default"
	awsRegion  = "ap-northeast-1"
	s3Bucket   = "pageloadtime-results"

	outPutFilePath = "../../analysis/pageload-status-code-info.csv"
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

// tokyo-01/example.com/example.com-with-cache-with-dane.csv
func targetCSV(s3ObjectKey string) (bool, string) {
	splitKey := strings.Split(s3ObjectKey, "/")
	if len(splitKey) != 3 {
		return false, ""
	}
	if strings.HasPrefix(splitKey[2], "letsdane-") || !strings.HasSuffix(splitKey[2], ".csv") {
		return false, ""
	}

	return true, splitKey[2]
}

type Records []Record

type Record struct {
	Status string
	Method string
	Domain string
}

func convertStruct(records [][]string) Records {
	var r Records
	if len(records) == 0 {
		return r
	}

	for _, record := range records[1:] {
		r = append(r, Record{
			Status: record[0],
			Method: record[1],
			Domain: record[2],
		})
	}
	return r
}

type Result []ResultRecord

type ResultRecord struct {
	MeasurementID string
	MeasurementInfo
	OneXX   int
	TwoXX   int
	ThreeXX int
	FourXX  int
	FiveXX  int
	XXX     int
	Total   int
}

func newResultRecord() *ResultRecord {
	return &ResultRecord{
		MeasurementID: "",
		MeasurementInfo: MeasurementInfo{
			Domain: "",
			Cache:  "",
			Dane:   "",
		},
		OneXX:   0,
		TwoXX:   0,
		ThreeXX: 0,
		FourXX:  0,
		FiveXX:  0,
		XXX:     0,
		Total:   0,
	}
}

type MeasurementInfo struct {
	Domain string
	Cache  string
	Dane   string
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

func (r *ResultRecord) setCalculatedResult(records Records) {
	for _, record := range records {
		statusCode, err := strconv.Atoi(record.Status)
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case statusCode >= 100 && statusCode < 200:
			r.OneXX++
		case statusCode >= 200 && statusCode < 300:
			r.TwoXX++
		case statusCode >= 300 && statusCode < 400:
			r.ThreeXX++
		case statusCode >= 400 && statusCode < 500:
			r.FourXX++
		case statusCode >= 500 && statusCode < 600:
			r.FiveXX++
		default:
			r.XXX++
		}
	}
	r.Total = r.OneXX + r.TwoXX + r.ThreeXX + r.FourXX + r.FiveXX + r.XXX
}

func exportResultAsCSV(result Result, filPath string) error {
	file, err := os.Create(filPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"measurementID", "domain", "cache", "dane", "1xx", "2xx", "3xx", "4xx", "5xx", "xxx", "total"}); err != nil {
		return err
	}

	for _, r := range result {
		csvRecord := []string{
			r.MeasurementID,
			r.Domain,
			r.Cache,
			r.Dane,
			strconv.Itoa(r.OneXX),
			strconv.Itoa(r.TwoXX),
			strconv.Itoa(r.ThreeXX),
			strconv.Itoa(r.FourXX),
			strconv.Itoa(r.FiveXX),
			strconv.Itoa(r.XXX),
			strconv.Itoa(r.Total),
		}
		if err := writer.Write(csvRecord); err != nil {
			return err
		}
	}
	return nil
}

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
	}

	// var wg sync.WaitGroup
	// var mutex sync.Mutex
	// sem := make(chan struct{}, 10)
	for idx, obj := range objs {
		logger.Info(fmt.Sprintf("now processing %s (%d/%d) %f%%", *obj.Key, idx+1, len(objs), float64(idx+1)/float64(len(objs))*100))
		ok, targetCSVFileName := targetCSV(*obj.Key)
		if !ok {
			continue
		}

		// wg.Add(1)
		// sem <- struct{}{}
		// go func(obj *s3.Object) {
		// defer func() {
		// 	<-sem
		// 	wg.Done()
		// }()

		gotObj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    obj.Key,
		})
		if err != nil {
			logger.Error(fmt.Sprintf("unable to get object %q from bucket %q, %v", *obj.Key, s3Bucket, err))
		}

		reader := csv.NewReader(gotObj.Body)
		records, err := reader.ReadAll()
		if err != nil {
			logger.Error(fmt.Sprintf("unable to read csv file %q, %v", *obj.Key, err))
		}

		convertedRecords := convertStruct(records)

		resultRecord := newResultRecord()
		resultRecord.setMeasurementID(*measurementID)
		resultRecord.setMeasurementInfoFromFile(targetCSVFileName)
		resultRecord.setCalculatedResult(convertedRecords)

		// mutex.Lock()
		result = append(result, *resultRecord)
		// mutex.Unlock()
		// }(obj)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].MeasurementID < result[j].MeasurementID
	})

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Domain < result[j].Domain
	})

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Cache < result[j].Cache
	})

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Dane < result[j].Dane
	})

	if err := exportResultAsCSV(result, *outPutFilePath); err != nil {
		logger.Error(fmt.Sprintf("unable to export result as csv, %v", err))
	}

	logger.Info(fmt.Sprintf("elapsed time: %f min", time.Since(start).Minutes()))
	logger.Info("done")
}
