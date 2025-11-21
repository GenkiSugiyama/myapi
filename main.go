package main

import (
	"log"
	"net/http"

	"github.com/GenkiSugiyama/myapi/handlers"
)

func main() {
	// ハンドラを宣言し変数に格納しただけではそのハンドラは使われない
	// サーバーで特定のハンドラを動かすための設定が必要になる
	// http.HandleFunc関数を使って、特定のパスに対してハンドラを紐付ける
	http.HandleFunc("/hello", handlers.HelloHandler)
	http.HandleFunc("/article", handlers.PostArticleHandler)
	http.HandleFunc("/article/list", handlers.ArticleListHandler)
	http.HandleFunc("/article/1", handlers.ArticleDetailHandler)
	http.HandleFunc("/article/nice", handlers.PostNiceHandler)
	http.HandleFunc("/comment", handlers.PostCommentHandler)

	log.Println("server start at port 8080")

	// http.ListenAndServe関数を使ってサーバーを起動する
	log.Fatal(http.ListenAndServe(":8080", nil))
}
