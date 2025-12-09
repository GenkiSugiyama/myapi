package controllers_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/GenkiSugiyama/myapi/controllers"
	"github.com/GenkiSugiyama/myapi/services"

	_ "github.com/go-sql-driver/mysql"
)

var aCon *controllers.ArticleController

func TestMain(m *testing.M) {
	dbUser := "user"
	dbPassword := "pass"
	dbDatabase := "myapi_db"
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Printf("DB setup failed: %v\n", err)
		os.Exit(1)
	}

	ser := services.NewMyAppService(db)
	aCon = controllers.NewArticleController(ser)

	m.Run()
}
