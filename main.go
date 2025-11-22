package main

import (
	"log"
	"net/http"

	"github.com/GenkiSugiyama/myapi/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// httpパッケージのデフォルトルーターではなく明示的にgorilla/muxのルータを宣言する
	r := mux.NewRouter()
	// ルータrのHandleFunc()でルーターにハンドラを登録する
	// mux.Route.Methods()で許可するHTTPメソッドを指定する
	r.HandleFunc("/hello", handlers.HelloHandler).Methods(http.MethodGet)
	r.HandleFunc("/article", handlers.PostArticleHandler).Methods(http.MethodPost)
	r.HandleFunc("/article/list", handlers.ArticleListHandler).Methods(http.MethodGet)
	// パスパラメータ{id}を含むルートを登録する
	r.HandleFunc("/article/{id:[0-9]+}", handlers.ArticleDetailHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/nice", handlers.PostNiceHandler).Methods(http.MethodPost)
	r.HandleFunc("/comment", handlers.PostCommentHandler).Methods(http.MethodPost)

	log.Println("server start at port 8080")

	// http.ListenAndServe（）の第二引数にはサーバーで使用するルーターを指定する必要がある
	// デフォルトのルーターを使用する場合はnilを指定する
	// 今回はgorilla/muxのルーターrを指定する
	log.Fatal(http.ListenAndServe(":8080", r))
}
