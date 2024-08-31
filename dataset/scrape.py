import argparse
import csv
import re
import time
import requests
from bs4 import BeautifulSoup



def scrape(url: str) -> BeautifulSoup:
    response = requests.get(url, timeout=30)
    html = response.text

    return BeautifulSoup(html, 'html.parser')

def scrape_from_file(filename: str) -> BeautifulSoup:
    with open(filename, 'r', encoding='utf-8') as file:
        return BeautifulSoup(file, 'html.parser')

def export_csv(result: list, filename: str) -> None:
    with open(filename, 'w') as file:
        writer = csv.writer(file)
        writer.writerows(result)

def main():
    parser = argparse.ArgumentParser(
        prog='scrape',
        description='Scrape the list of websites from the given HTML file. If the input HTML file is not provided, this program will scrape the website(https://internet.nl/halloffame/web/) directly.'
    )
    parser.add_argument('-in', '--input_html', help="Path to the HTML file to be parsed.")
    parser.add_argument('-out', '--output_csv', help="Path to the CSV file to be exported.")
    args = parser.parse_args()

    print("Current time: ", time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
    if args.input_html:
        soup = scrape_from_file(args.input_html)
    else:
        url = "https://internet.nl/halloffame/web/"
        soup = scrape(url)

    result = []
    ul_element = soup.find('ul', class_='list-column column-3')
    if ul_element:
        li_elements = ul_element.find_all('li')
        for idx, li_element in enumerate(li_elements):
            a_element = li_element.find('a')
            if a_element:
                href = a_element.get('href')

            # https://internet.nl/site/edienstenburgerzaken-test.barendrecht.nl/2371223/ -> edienstenburgerzaken-test.barendrecht.nl
            match = re.search(r'/site/([^/]+)/', href)
            if match:
                domain = match.group(1)
                result.append([domain])

    export_csv(result, args.output_csv)

if __name__ == '__main__':
    main()
