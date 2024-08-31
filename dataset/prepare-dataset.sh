
# 1. save html file of hall of fame websites
html=$(python3 save-html.py)


# 2. scrape all domains from the html file
python3 scrape.py -in=${html} -out=hall-of-flame-websites.csv


# 3. extract TLSA records for each domain
cd ../cmd/zdns/tlsa-record/
go run main.go -inputCSV ../../../dataset/hall-of-flame-websites.csv -outputCSV ../../../dataset/hall-of-flame-websites-tlsa.csv
cd "$OLDPWD"

# 4. extract TLSA recotds whose usage is 3
python3 usage-3.py -in=hall-of-flame-websites-tlsa.csv -out=hall-of-flame-websites-tlsa-usage-3.csv
