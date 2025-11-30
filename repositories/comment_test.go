package repositories_test

import (
	"testing"

	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/repositories"
)

func TestFindArticleCommentsByArticleID(t *testing.T) {
	articleID := 1
	got, err := repositories.FindArticleCommentsByArticleID(testDB, articleID)
	if err != nil {
		t.Fatal(err)
	}

	for _, comment := range got {
		if comment.ArticleID != articleID {
			t.Errorf("want comment of articleID %d but got ID %d\n", articleID, comment.ArticleID)
		}
	}
}

func TestInsertComment(t *testing.T) {
	comment := models.Comment{
		ArticleID: 1,
		Message:   "comment insert test",
	}

	expectedCommentID := 3
	newComment, err := repositories.InsertComment(testDB, comment)
	if err != nil {
		t.Fatal(err)
	}

	if newComment.CommentID != expectedCommentID {
		t.Errorf("new comment id is expected %d but got %d\n", expectedCommentID, newComment.CommentID)
	}

	t.Cleanup(func() {
		const sqlStr = `
			DELETE FROM comments
			WHERE message = ?;
		`
		testDB.Exec(sqlStr, comment.Message)
	})
}
