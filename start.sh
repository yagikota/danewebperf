#!/bin/bash

measurementID=$1
last=$2
concurrency=$3

cd cmd/pageloadtime
go build -o pageloadtime

# this is dataset
inputCSV=../../dataset/hall-of-flame/hall-of-flame-websites-tlsa-usage3.csv

mkdir -p ../../result/pageloadtime/${measurementID}

echo "Running measurements without cache and without DANE...(1/4)"
./pageloadtime -subdirname=${measurementID} -inputCSV=${inputCSV} -last=${last} -concurrency=${concurrency} > ../../result/pageloadtime/${measurementID}/without-cache-without-dane.log
echo "Uplodaing results to S3...(1/4)"
aws s3 mv ../../result/pageloadtime/${measurementID}/ s3://pageloadtime-results/${measurementID}/ --recursive

echo "Running measurements without cache and with DANE...(2/4)"
./pageloadtime -subdirname=${measurementID} -inputCSV=${inputCSV} -last=${last}  -concurrency=${concurrency} -dane > ../../result/pageloadtime/${measurementID}/without-cache-with-dane.log
echo "Uplodaing results to S3...(2/4)"
aws s3 mv ../../result/pageloadtime/${measurementID}/ s3://pageloadtime-results/${measurementID}/ --recursive

echo "Running measurements with cache and without DANE...(3/4)"
./pageloadtime -subdirname=${measurementID} -inputCSV=${inputCSV} -last=${last} -concurrency=${concurrency} -cache > ../../result/pageloadtime/${measurementID}/with-cache-without-dane.log
echo "Uplodaing results to S3...(3/4)"
aws s3 mv ../../result/pageloadtime/${measurementID}/ s3://pageloadtime-results/${measurementID}/ --recursive

echo "Running measurements with cache and with DANE... (4/4)"
./pageloadtime -subdirname=${measurementID} -inputCSV=${inputCSV} -last=${last} -concurrency=${concurrency} -cache -dane > ../../result/pageloadtime/${measurementID}/with-cache-with-dane.log
echo "Uplodaing results to S3...(4/4)"
aws s3 mv ../../result/pageloadtime/${measurementID}/ s3://pageloadtime-results/${measurementID}/ --recursive
