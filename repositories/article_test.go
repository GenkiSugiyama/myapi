package repositories_test

import (
	"testing"

	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/repositories"
	"github.com/GenkiSugiyama/myapi/repositories/testdata"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetArticleDetailByID(t *testing.T) {
	// テーブルドリブンテスト
	tests := []struct {
		testTitle string
		expected  models.Article
	}{
		{
			testTitle: "subtest1",
			expected:  testdata.ArticleTestData[0],
		},
		{
			testTitle: "subtest2",
			expected:  testdata.ArticleTestData[1],
		},
	}

	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			got, err := repositories.GetArticleDetailByID(testDB, test.expected.ID)
			if err != nil {
				t.Fatal(err)
			}

			if got.ID != test.expected.ID {
				t.Errorf("ID: got %d but want %d\n", got.ID, test.expected.ID)
			}
			if got.Title != test.expected.Title {
				t.Errorf("Title: got %s but want %s\n", got.Title, test.expected.Title)
			}
			if got.Contents != test.expected.Contents {
				t.Errorf("Contents: got %s but want %s\n", got.Contents, test.expected.Contents)
			}
			if got.UserName != test.expected.UserName {
				t.Errorf("UserName: got %s but want %s\n", got.UserName, test.expected.UserName)
			}
			if got.NiceNum != test.expected.NiceNum {
				t.Errorf("NiceNum: got %d but want %d\n", got.NiceNum, test.expected.NiceNum)
			}
		})
	}
}

func TestFindArticles(t *testing.T) {
	expectedNum := len(testdata.ArticleTestData)
	got, err := repositories.FindArticles(testDB, 1)
	if err != nil {
		t.Fatal(err)
	}
	if num := len(got); num != expectedNum {
		t.Errorf("expected %d articles, got %d", expectedNum, num)
	}
}

func TestInsertArticle(t *testing.T) {
	article := models.Article{
		Title:    "testTitle",
		Contents: "testContents",
		UserName: "testUser",
	}

	expectedArticleNum := 3
	newArticle, err := repositories.InsertArticle(testDB, article)
	if err != nil {
		t.Error(err)
	}
	if newArticle.ID != expectedArticleNum {
		t.Errorf("new article id is expected %d, got %d\n", expectedArticleNum, newArticle.ID)
	}

	t.Cleanup(func() {
		const sqlStr = `
			DELETE FROM articles 
			WHERE title = ? AND contents = ? AND username = ?;
		`
		testDB.Exec(sqlStr, article.Title, article.Contents, article.UserName)
	})
}

func TestUpdateArticleNice(t *testing.T) {
	articleID := 1
	before, err := repositories.GetArticleDetailByID(testDB, articleID)
	if err != nil {
		t.Fatal("fail to get before data")
	}

	err = repositories.UpdateArticleNice(testDB, articleID)
	if err != nil {
		t.Fatal(err)
	}

	after, err := repositories.GetArticleDetailByID(testDB, articleID)
	if err != nil {
		t.Fatal("fail to get after data")
	}

	if after.NiceNum-before.NiceNum != 1 {
		t.Error("fail to update nice num")
	}

	t.Cleanup(func() {
		const sqlStr = `
			UPDATE articles
			SET nice = ?
			WHERE article_id = ?;
		`
		testDB.Exec(sqlStr, before.NiceNum, articleID)
	})
}
