# Data Set

- `hall-of-flame-websites-2023-10-27-00-14-00.html.zip`
  - This is a zip file of the HTML file of this site(<https://internet.nl/halloffame/>), which downloaded at 2023-10-27 00:14:00. This is the base of data set to used in this research. This site is updated every day, and the data set is not stable. So we use the data set that is downloaded in advance.

- `hall_of_flame_websites.csv`
  - This is a csv file of **31696** websites that are listed in this site(<https://internet.nl/halloffame/>).

- `hall-of-flame-websites-tlsa.csv`
  - This is a csv file of websites that are listed in this site(<https://internet.nl/halloffame/>), and have TLSA record. These websites are looked up by using [zdns](https://github.com/zmap/zdns)

- `hall-of-flame-websites-tlsa-usage3.csv`
  - This is a csv file of **3652** websites that are listed in this site(<https://internet.nl/halloffame/>), and have TLSA record, and have TLSA record usage 3.

- `nu-tlsa-one-mx-2023-05-24.txt`
  - <https://data.openintel.nl/data/open-tld/2023/tlsa/nu-tlsa-one-mx-2023-05-24.txt>
- `se-tlsa-one-mx-2023-05-24.txt`
  - <https://data.openintel.nl/data/open-tld/2023/tlsa/se-tlsa-one-mx-2023-05-24.txt>
In this research, we use `hall-of-flame-websites-tlsa-usage3.csv` as data set, because we want to analyze the websites that use TLSA record usage 3 that is the most common.
