package repositories

import (
	"database/sql"
	"fmt"

	"github.com/GenkiSugiyama/myapi/models"
)

func InsertArticle(db *sql.DB, article models.Article) (models.Article, error) {
	const sqlStr = `
		INSERT INTO articles (title, contents, username, nice, created_at) 
		VALUES (?, ?, ?, 0, NOW());
		`

	newArticle := models.Article{
		Title:    article.Title,
		Contents: article.Contents,
		UserName: article.UserName,
	}

	result, err := db.Exec(sqlStr, article.Title, article.Contents, article.UserName)
	if err != nil {
		fmt.Println(err)
		return models.Article{}, err
	}

	newArticleID, _ := result.LastInsertId()
	newArticle.ID = int(newArticleID)

	return newArticle, nil
}
