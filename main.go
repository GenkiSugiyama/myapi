package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GenkiSugiyama/myapi/api"
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

	// controller層の構造体の具体的な生成処理をルータ層に移したことによってmain.goがアプリケーションの起動に必要な処理だけになった
	r := api.NewRouter(db)

	log.Println("server start at port 8080")

	// http.ListenAndServe（）の第二引数にはサーバーで使用するルーターを指定する必要がある
	// デフォルトのルーターを使用する場合はnilを指定する
	// 今回はgorilla/muxのルーターrを指定する
	log.Fatal(http.ListenAndServe(":8080", r))
}
