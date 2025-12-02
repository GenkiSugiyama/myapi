package repositories

import (
	"database/sql"

	"github.com/GenkiSugiyama/myapi/models"
)

func FindArticleCommentsByArticleID(db *sql.DB, articleID int) ([]models.Comment, error) {
	const sqlStr = `
		SELECT comment_id, article_id, message, created_at
		FROM comments
		WHERE article_id = ?
	`

	rows, err := db.Query(sqlStr, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		var createdAt sql.NullTime
		err := rows.Scan(&comment.CommentID, &comment.ArticleID, &comment.Message, &createdAt)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			comment.CreatedAt = createdAt.Time
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func InsertComment(db *sql.DB, comment models.Comment) (models.Comment, error) {
	const sqlStr = `
		INSERT INTO comments (article_id, message, created_at)
		VALUES (?, ?, NOW())
	`

	newComment := models.Comment{
		ArticleID: comment.ArticleID,
		Message:   comment.Message,
	}

	result, err := db.Exec(sqlStr, newComment.ArticleID, newComment.Message)
	if err != nil {
		return models.Comment{}, err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return models.Comment{}, err
	}

	newComment.CommentID = int(newID)

	return newComment, nil
}
