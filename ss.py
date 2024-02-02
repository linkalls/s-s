import os
import requests
from bs4 import BeautifulSoup
import threading
import PySimpleGUI as sg

def download_images(url, folder_name, progress_bar):
    # URLからHTMLを取得
    response = requests.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')
    
    # タイトルを抽出
    title = soup.find('a', itemprop='url').text
    
    # フォルダ名を作成
    folder_name = os.path.join(folder_name, title)
    
    # フォルダが存在しない場合は作成
    if not os.path.exists(folder_name):
        os.makedirs(folder_name)
    
    # 画像要素を抽出
    images = soup.find_all('img', {'class': 'pict'})
    
    # 画像をダウンロードして保存
    for i, img in enumerate(images):
        image_url = img['src'].replace('-s.jpg', '.jpg')
        image_data = requests.get(image_url).content
        with open(os.path.join(folder_name, f'{i}.jpg'), 'wb') as handler:
            handler.write(image_data)

        # 進行状況バーの更新
        progress_bar.UpdateBar(current_count=i + 1, max=len(images))

def start_download(window, values):
    url = values['-IN-']
    folder_name = values['-FOLDER-']
    # 画像の数を取得して進行状況バーを初期化
    response = requests.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')
    max_images = len(soup.find_all('img', {'class': 'pict'}))
    window['-PROGRESS-'].UpdateBar(current_count=0, max=max_images)
    threading.Thread(target=download_images, args=(url, folder_name, window['-PROGRESS-'])).start()

# GUIの初期化
layout = [[sg.Text("Enter URL:")],
          [sg.Input(key='-IN-')],
          [sg.Text("Enter Folder Path:")],
          [sg.Input(key='-FOLDER-')],
          [sg.Button("Start Download"), sg.Button("Cancel")],
          [sg.ProgressBar(max_value=1, size=(20, 20), key='-PROGRESS-')]]

window = sg.Window('Image Downloader', layout)

while True:
    event, values = window.read()
    if event == "Start Download":
        start_download(window, values)
    elif event == "Cancel" or event == sg.WINDOW_CLOSED:
        break

window.close()