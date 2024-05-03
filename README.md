# DANEWebPerf

DANEWebPerf is a tool to measure the performance of DANE-enabled websites.

## Installation

``` shell
git clone --recursive git@github.com:yagikota/danewebperf.git
```


### prepare python virtual environment

```bash
python3 -m venv python-env
source ./python-env/bin/activate
```

## Usage

### Build Docker Images

``` shell
./prepare.sh
```

### Run the Measurement

``` shell
./start.sh
```

## How to Analyze the Measurement Results

### copy the pageload time csv files to the `result/pageloadtime/` directory

``` shell
cd result/pageloadtime/
./s3-pageload-csv-copy.sh
```
then, you can find the pageload time csv files like this:

``` shell
tree
.
├── frankfurt-01
│   ├── pageloadtime-with-cache-with-dane.csv
│   ├── pageloadtime-with-cache-without-dane.csv
│   ├── pageloadtime-without-cache-with-dane.csv
│   └── pageloadtime-without-cache-without-dane.csv
├── frankfurt-02
│   ├── pageloadtime-with-cache-with-dane.csv
│   ├── pageloadtime-with-cache-without-dane.csv
│   ├── pageloadtime-without-cache-with-dane.csv
│   └── pageloadtime-without-cache-without-dane.csv
├── s3-pageload-csv-copy.sh

cat frankfurt-01/pageloadtime-with-cache-with-dane.csv | head
domain,pageLoadTime,cache,dane
06sex.nl,224,true,true
0x1a8510f2.space,,true,true
112-ov.nl,1823,true,true
123derivatives.com,,true,true
158.nl,384,true,true
19hetatelier.nl,,true,true
1dos.eu,,true,true
1fc0.de,,true,true
200mmx.net,160,true,true
```
