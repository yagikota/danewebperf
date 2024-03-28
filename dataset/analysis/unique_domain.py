import pandas as pd
import argparse

def print_unique_domain(df):
    unique_domain_num = len(df['domain'].unique())
    print(unique_domain_num)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        prog='unique_domain',
        description='Print unique domain'
    )
    parser.add_argument('-in', '--input_csv', type=str, help='The input csv file path')
    args = parser.parse_args()

    try:
        df = pd.read_csv(args.input_csv)
        print_unique_domain(df)
    except Exception as e:
        print(e)
