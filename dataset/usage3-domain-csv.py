import pandas as pd
import argparse


def unique_index_name_by_usage(df, usage):
    filtered_new_tls_df = df[df['usage'] == usage]
    unique_index_name = filtered_new_tls_df[['domain']].drop_duplicates()
    sorted_unique_index_name = unique_index_name.sort_values(by=['domain'])
    return sorted_unique_index_name

def export_domain_by_usage_to_csv(df, csv_path):
    df.to_csv(csv_path, index=False, header=False)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        prog='usage3_domain_csv',
        description='Export domain that has TLSA record with usage 3'
    )
    parser.add_argument('-in', '--input_csv', type=str, help='The input csv file path')
    parser.add_argument('-out', '--output_csv', type=str, help='The output csv file path')
    args = parser.parse_args()

    try:
        csv_path = args.input_csv
        df = pd.read_csv(csv_path)
        df = unique_index_name_by_usage(df, 3)
        export_domain_by_usage_to_csv(df, args.output_csv)
    except Exception as e:
        print(e)
