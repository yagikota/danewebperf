# this program read given csv files that listing domain name, and merge them into one csv file with unique domain name
#

import pandas as pd
import argparse


def merge_csv(csv_list, output_csv):
    df_list = []
    for csv in csv_list:
        df_list.append(pd.read_csv(csv, names=['domain']))
    df = pd.concat(df_list, ignore_index=True)
    df = df.drop_duplicates()
    df = df.sort_values(by=['domain'])
    df.to_csv(output_csv, index=False, header=False)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        prog='merge_dataset',
        description='Merge multiple csv files into one csv file with unique domain name'
    )
    parser.add_argument('-in', '--input_csv', nargs='+', type=str, help='The input csv file path')
    parser.add_argument('-out', '--output_csv', type=str, help='The output csv file path')
    args = parser.parse_args()

    try:
        merge_csv(args.input_csv, args.output_csv)
    except Exception as e:
        print(e)
