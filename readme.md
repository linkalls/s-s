

# プロジェクトの使用方法

このプロジェクトは、指定されたURLから画像をダウンロードし、それらの画像をPDFに変換するためのツールです。以下の手順で使用できます。

## 1. 画像のダウンロード

まず、`@yuyu.go`を使用して画像をダウンロードします。以下のコマンドを実行してください。

```bash
go run yuyu.go -txt path/to/your/textfile.txt
```

または、直接URLを入力して画像をダウンロードすることもできます。

```bash
go run yuyu.go https://example.com/archives/example
```

このコマンドを実行すると、`@yuyu.go`は指定されたテキストファイルからURLを読み取り、それらのURLから画像をダウンロードします。ダウンロードされた画像は、`USERPROFILE/Downloads/yuyublog`ディレクトリ内の各フォルダに保存されます。

## 2. 画像のPDF変換

ダウンロードした画像をPDFに変換するには、`@pngtopdf.go`を使用します。以下のコマンドを実行してください。

```bash
go run pngtopdf.go
```

このコマンドを実行すると、`@pngtopdf.go`は`USERPROFILE/Downloads/yuyublog`ディレクトリ内の各フォルダをスキャンし、それぞれのフォルダに含まれる画像をPDFに変換します。変換されたPDFは、`USERPROFILE/Downloads/yuyublog/pdf`ディレクトリに保存されます。

## 注意事項

- ダウンロードされた画像のフォルダ名には、Windowsのファイル名に使用できない文字が含まれている場合があります。これらの文字は`@yuyu.go`によって削除または置換されます。
- ダウンロードされた画像のファイル名は、`@yuyu.go`によって生成され、`@pngtopdf.go`でPDFに変換する際に使用されます。

---

このREADMEは、プロジェクトの基本的な使用方法を説明しています。プロジェクトの詳細な機能やオプションについては、各ファイルのコードを確認してください。