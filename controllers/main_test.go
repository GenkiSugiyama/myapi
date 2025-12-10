package controllers_test

import (
	"testing"

	"github.com/GenkiSugiyama/myapi/controllers"
	"github.com/GenkiSugiyama/myapi/controllers/testdata"

	_ "github.com/go-sql-driver/mysql"
)

var aCon *controllers.ArticleController

func TestMain(m *testing.M) {
	ser := testdata.NewServiceMock()
	aCon = controllers.NewArticleController(ser)

	m.Run()
}
