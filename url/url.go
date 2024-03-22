package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	baseUrl := "https://yuyubox.blog.jp/"
	downloadDir := "./" // 現在のディレクトリに保存する場合

	// url.txtファイルを開く
	filePath := filepath.Join(downloadDir, "url.txt")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("エラー: ファイルの作成に失敗しました: %s\n", filePath)
		return
	}
	defer f.Close()

	// 既存のURLをセットに保存
	urlSet := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urlSet[scanner.Text()] = true
	}

	page := 1
	for {
		url := fmt.Sprintf("%s?p=%d", baseUrl, page)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("エラー: URLの取得に失敗しました: %s\n", url)
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println("エラー: goqueryドキュメントの作成に失敗しました")
			return
		}

		var found bool
		doc.Find("h1.article-title a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && !urlSet[href] {
				found = true
				// URLをファイルに保存
				_, err := f.WriteString(href + "\n")
				if err != nil {
					fmt.Printf("エラー: URLの保存に失敗しました: %s\n", href)
				}
			}
		})

		if !found {
			break
		}

		page++
	}

	fmt.Println("URLの抽出が完了しました。")
}