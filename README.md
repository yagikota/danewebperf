# danewebperf

DANEWebPerf is a tool for measuring page load time of websites with DANE.

## Create Python Virtual Environment

```bash
python3 -m venv python-env
source ./python-env/bin/activate
```

## Prepare DataSet

``` shell
python3 save-html.py
```

``` shell
python3 scrape.py --input_html=hall-of-flame-websites-2024-03-06-17-18-17.html --output_csv=hall-of-flame-websites.csv
```

``` shell
cd cmd/zdns/tlsa-record
nohup go run main.go &
```

## Prepare ZDNS

``` shell
git submodule add https://github.com/zmap/zdns.git  zdns
cd zdns
go build
mv zdns  ../cmd/zdns/
```
