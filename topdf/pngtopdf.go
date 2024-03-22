package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// yuyu.goで保存した画像があるフォルダのパスを指定
	downloadDir := filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "yuyublog")

	// PDFフォルダが存在しない場合は作成
	pdfDir := filepath.Join(downloadDir, "pdf")
	if _, err := os.Stat(pdfDir); os.IsNotExist(err) {
		os.Mkdir(pdfDir, os.ModePerm)
	}

	// フォルダ内の各フォルダを処理
	err := filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// フォルダ内のフォルダを処理
		if info.IsDir() {
			// PDFを作成
			pdf := gofpdf.New("P", "mm", "A4", "")

			// フォルダ内の画像ファイルを処理
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				// 新しいページを作成し、画像を追加
				pdf.AddPage()
				pdf.Image(filepath.Join(path, file.Name()), 10, 10, 190, 0, false, "", 0, "")
			}

			// PDFファイル名を生成
			pdfName := info.Name() + ".pdf"
			pdfPath := filepath.Join(pdfDir, pdfName)

			// PDFを保存
			err = pdf.OutputFileAndClose(pdfPath)
			if err != nil {
				fmt.Printf("エラー: PDFの保存に失敗しました: %s\n", pdfPath)
			} else {
				fmt.Printf("画像をPDFに変換しました: %s\n", pdfName)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("エラー: フォルダの処理に失敗しました")
	}
}