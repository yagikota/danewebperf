# danewebperf

DANEWebPerf is a tool to measure the performance of DANE-enabled websites.

## Requirements

### prepare python virtual environment

```bash
python3 -m venv python-env
source ./python-env/bin/activate
```

### build ZDNS

``` shell
git submodule add https://github.com/zmap/zdns.git  zdns
cd zdns
go build
mv zdns  ../cmd/zdns/
```

### prepare Data Set

``` shell
prepare-data-set.sh
```

## Usage

### 
