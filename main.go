package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GenkiSugiyama/myapi/controllers"
	"github.com/GenkiSugiyama/myapi/services"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbDatabase = os.Getenv("DB_NAME")
	dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
)

func main() {
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Printf("failed to connect DB: %v\n", err)
	}

	ser := services.NewMyAppService(db)
	con := controllers.NewMyAppController(ser)
	// httpパッケージのデフォルトルーターではなく明示的にgorilla/muxのルータを宣言する
	r := mux.NewRouter()
	// ルータrのHandleFunc()でルーターにハンドラを登録する
	// mux.Route.Methods()で許可するHTTPメソッドを指定する
	r.HandleFunc("/hello", con.HelloHandler).Methods(http.MethodGet)
	r.HandleFunc("/article", con.PostArticleHandler).Methods(http.MethodPost)
	r.HandleFunc("/article/list", con.ArticleListHandler).Methods(http.MethodGet)
	// パスパラメータ{id}を含むルートを登録する
	r.HandleFunc("/article/{id:[0-9]+}", con.ArticleDetailHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/nice", con.PostNiceHandler).Methods(http.MethodPost)
	r.HandleFunc("/comment", con.PostCommentHandler).Methods(http.MethodPost)

	log.Println("server start at port 8080")

	// http.ListenAndServe（）の第二引数にはサーバーで使用するルーターを指定する必要がある
	// デフォルトのルーターを使用する場合はnilを指定する
	// 今回はgorilla/muxのルーターrを指定する
	log.Fatal(http.ListenAndServe(":8080", r))
}
