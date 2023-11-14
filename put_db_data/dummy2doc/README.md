# ダミーデータ投入ツール

## 事前準備

1. 環境変数の設定
	- `.env.example` ファイルから `.env` ファイルを作成
		```
		cp .env.example .env
		```
	- `.env` にDBの接続情報を追記
		```text:.env
		ENDPOINT=example.com
		PORT=3306
		DBNAME=example
		USERNAME=example
		PASSWORD=example
		```
1. 仮想環境を作成
	- 仮想環境を作成するディレクトリに移動
		```
		cd tools/put_db_data/dummy2doc
		```
	- 仮想環境を作成
		```
		python -m venv venv
		```
	- 確認
		```
		ls 
		```
		- カレントディレクトリに `venv` ファイルが作成されていること
1. 仮想環境へ入る
	```
	source venv/bin/activate
	```
	- コマンドラインに表示されるディレクトリの先頭に `(venv)` がついていたら成功
		```
		// 例
		(venv) root:~/github.com/tsuruya-jp/tools/put_db_data/dummy2doc
		```
1. ライブラリをインストール
	```
	pip install -r requirements.txt
	```

## 使用方法
1. 投入したいダミーデータをcsvファイルに書く
	- データのカラムは以下の通り
		```
		"オーダー番号","発行日","氏名","金額","端末番号","ユーザー番号","会社名"
		```
1. 対象のcsvファイルを `put_db_data/dummy2doc` ディレクトリへ置く
1. `python` コマンドでダミーツールを実行
	- 実行には引数が必要です
	```
	required:	invoice or receipt or quotation
	optional:	-h, --help	"show this help message and exit"	
			-f, --file	"optional: Please select input data file"
				default: dummy.csv
	```
	- (例) python exec.py invoice

## 実行後
1. 仮想環境から出る
	```
	deactivate
	```
