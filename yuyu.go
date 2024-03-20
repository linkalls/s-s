package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	txtPath := flag.String("txt", "", "テキストファイルのパス")
	flag.Parse()

	// USERPROFILEの/Downloads内にyuyublogを作成
	downloadDir := filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "yuyublog")
	os.MkdirAll(downloadDir, os.ModePerm)

	if *txtPath == "" {
		fmt.Print("URLを入力してください: ")
		var urlStr string
		fmt.Scanln(&urlStr)
		DownloadImages(urlStr, downloadDir)
	} else {
		file, err := os.Open(*txtPath)
		if err != nil {
			fmt.Println("エラー: ファイルのオープンに失敗しました")
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			DownloadImages(scanner.Text(), downloadDir)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("エラー: ファイルの読み取りに失敗しました")
		}
	}
}

func DownloadImages(urlStr string, downloadDir string) {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("エラー: URLの取得に失敗しました")
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("エラー: goqueryドキュメントの作成に失敗しました")
		return
	}

	// <h1 class="article-title">タグの中の<a>タグからタイトルを取得
	title := doc.Find("h1.article-title a").First().Text()
	// タイトルをフォルダ名として使用して、そのフォルダを作成
	folderPath := filepath.Join(downloadDir, sanitizeFileName(title))
	os.MkdirAll(folderPath, os.ModePerm)

		imageIndex := 1
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			imageUrl, exists := s.Attr("href")
			if exists && isImageUrl(imageUrl) {
				err := downloadImage(imageUrl, folderPath, imageIndex)
				if err != nil {
					fmt.Printf("エラー: 画像のダウンロードに失敗しました: %s\n", imageUrl)
				} else {
					imageIndex++
				}
			}
		})
	}

func isImageUrl(url string) bool {
	// 画像URLを指しているかどうかを判定する関数
	// ここでは、URLが.jpg、.jpeg、.png、.gif、.webpの拡張子を持っているかどうかで判定しています
	return filepath.Ext(url) == ".jpg" || filepath.Ext(url) == ".jpeg" || filepath.Ext(url) == ".png" || filepath.Ext(url) == ".gif" || filepath.Ext(url) == ".webp"
}

func downloadImage(imageUrl string, downloadDir string, index int) error {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// ファイル名を生成
	fileName := fmt.Sprintf("%03d.jpg", index)
	filePath := filepath.Join(downloadDir, fileName)

	// ファイルを保存
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("画像をダウンロードしています: %s\n", fileName)
	return nil
}

func sanitizeFileName(input string) string {
	// Windowsのファイル名に使用できない文字を削除または置換します
	// また、.webp拡張子を.jpgに置換します
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := re.ReplaceAllString(input, "")
	// .webp拡張子を.jpgに置換します
	sanitized = strings.Replace(sanitized, ".webp", ".jpg", -1)
	return sanitized
}