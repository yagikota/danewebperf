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
