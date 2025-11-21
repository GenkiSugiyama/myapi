package handlers

import (
	"fmt"
	"io"
	"net/http"
)

// リクエストを受け取って任意のレスポンスを書き込むための関数型ハンドラを宣言する
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// どのようなレスポンスを返すのかを記述する
	// http.ResponseWriterに対してHello, world!と書き込む
	io.WriteString(w, "Hello, world!\n")
}

func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Posting Article...\n")
}

func ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Article List\n")
}

func ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	articleID := 1
	resString := fmt.Sprintf("Article Detail: No.%d\n", articleID)
	io.WriteString(w, resString)
}

func PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Posting Nice...\n")
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Posting Comment...\n")
}
