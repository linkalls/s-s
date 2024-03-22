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

		  // URLを保存するスライス
			var urls []string

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
					urls = append(urls, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
					fmt.Println("エラー: ファイルの読み取りに失敗しました")
			}

			// URLを逆順に処理
			for i := len(urls) - 1; i >= 0; i-- {
					DownloadImages(urls[i], downloadDir)
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

    // タイトルの末尾にある全角数字を [3] のような形式に変更します
    re = regexp.MustCompile(`（(\d+)）$`)
    sanitized = re.ReplaceAllString(sanitized, "[$1]")

    // 全角数字表記を一般的な数字に変換します
    re = regexp.MustCompile(`①|②|③|④|⑤|⑥|⑦|⑧|⑨|⑩|⑪|⑫|⑬|⑭|⑮|⑯|⑰|⑱|⑲`)
    sanitized = re.ReplaceAllStringFunc(sanitized, func(s string) string {
        switch s {
        case "①":
            return "[1]"
        case "②":
            return "[2]"
        case "③":
            return "[3]"
        case "④":
            return "[4]"
        case "⑤":
            return "[5]"
        case "⑥":
            return "[6]"
        case "⑦":
            return "[7]"
        case "⑧":
            return "[8]"
        case "⑨":
            return "[9]"
        case "⑩":
            return "[10]"
        case "⑪":
            return "[11]"
        case "⑫":
            return "[12]"
        case "⑬":
            return "[13]"
        case "⑭":
            return "[14]"
        case "⑮":
            return "[15]"
        case "⑯":
            return "[16]"
        case "⑰":
            return "[17]"
        case "⑱":
            return "[18]"
        case "⑲":
            return "[19]"
        default:
            return s
        }
    })

    return sanitized
}