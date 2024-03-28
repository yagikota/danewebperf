
import pandas as pd
import argparse

def count_domains_by_tld(csv):
    df = pd.read_csv(csv, names=['domain'])
    tlds = df['domain'].str.split('.').str[-1]
    tld_count = tlds.value_counts()
    print(tld_count)

if __name__ == '__main__':
    pd.set_option('display.max_rows', None)
    parser = argparse.ArgumentParser(
        prog='count_domains_by_tld',
        description='Count domain by tld'
    )
    parser.add_argument('-in', '--input_csv', type=str, help='The input csv file path')
    args = parser.parse_args()

    try:
        count_domains_by_tld(args.input_csv)
    except Exception as e:
        print(e)

