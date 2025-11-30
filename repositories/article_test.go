package repositories_test

import (
	"testing"

	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/repositories"
)

func TestGetArticleDetailByID(t *testing.T) {
	dbUser := "user"
	dbPassword := "pass"
	dbDatabase := "myapi_db"
	expected := models.Article{}

	got, err := repositories.GetArticleDetailByID()
	if err != nil {
		t.Fatal(err)
	}

	if got != expected {
		t.Errorf("got %s but want %s\n", got, expected)
	}
}
