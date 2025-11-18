package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// リクエストを受け取って任意のレスポンスを書き込むための関数型ハンドラを宣言する
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		// どのようなレスポンスを返すのかを記述する
		// http.ResponseWriterに対してHello, world!と書き込む
		io.WriteString(w, "Hello, world!\n")
	}

	// ハンドラを宣言し変数に格納しただけではそのハンドラは使われない
	// サーバーで特定のハンドラを動かすための設定が必要になる
	// http.HandleFunc関数を使って、特定のパスに対してハンドラを紐付ける
	http.HandleFunc("/", helloHandler)

	log.Println("server start at port 8080")

	// http.ListenAndServe関数を使ってサーバーを起動する
	log.Fatal(http.ListenAndServe(":8080", nil))
}
