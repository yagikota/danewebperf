import requests
import datetime

url = "https://internet.nl/halloffame/web/"

# URLからHTMLを取得
response = requests.get(url)

# ステータスコードが200（成功）の場合のみ処理を続行
if response.status_code == 200:
    current_time = datetime.datetime.now()
    html_filename = f"hall-of-flame-websites-{current_time.year}-{current_time.month:02d}-{current_time.day:02d}-{current_time.hour:02d}-{current_time.minute:02d}-{current_time.second:02d}.html"
    with open(html_filename, "w", encoding="utf-8") as f:
        f.write(response.text)