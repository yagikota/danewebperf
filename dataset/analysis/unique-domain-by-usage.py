import pandas as pd
import argparse

def print_unique_domain_by_usage(df, usage):
    filtered_new_tls_df = df[df['usage'] == usage]
    unique_index_count_new_tls = len(filtered_new_tls_df['domain'].unique())
    print(unique_index_count_new_tls)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        prog='unique_domain_by_usage',
        description='Print unique domain that has TLSA record with usage'
    )
    parser.add_argument('-in', '--input_csv', type=str, help='The input csv file path')
    args = parser.parse_args()

    try:
        df = pd.read_csv(args.input_csv)
        print_unique_domain_by_usage(df, 0)
        print_unique_domain_by_usage(df, 1)
        print_unique_domain_by_usage(df, 2)
        print_unique_domain_by_usage(df, 3)
    except Exception as e:
        print(e)
