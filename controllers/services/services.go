package services

import "github.com/GenkiSugiyama/myapi/models"

// controllerが必要とするサービスのインターフェースをコントローラ層で定義する
type ArticleServicer interface {
	PostArticleService(article models.Article) (models.Article, error)
	ArticleListService(page int) ([]models.Article, error)
	GetArticleService(articleID int) (models.Article, error)
	PostNiceService(article models.Article) (models.Article, error)
}

type CommentServicer interface {
	PostCommentService(comment models.Comment) (models.Comment, error)
}
