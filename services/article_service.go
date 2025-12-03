package services

import (
	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/repositories"
)

func PostArticleService(article models.Article) (models.Article, error) {
	db, err := connectDB()
	if err != nil {
		return models.Article{}, err
	}
	defer db.Close()

	newArticle, err := repositories.InsertArticle(db, article)
	if err != nil {
		return models.Article{}, err
	}

	return newArticle, nil
}

func ArticleListService(page int) ([]models.Article, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	articleList, err := repositories.FindArticles(db, page)
	if err != nil {
		return nil, err
	}

	return articleList, nil
}

func GetArticleService(articleID int) (models.Article, error) {
	// sql.DB型を取得する
	db, err := connectDB()
	if err != nil {
		return models.Article{}, err
	}
	defer db.Close()

	article, err := repositories.GetArticleDetailByID(db, articleID)
	if err != nil {
		return models.Article{}, err
	}

	commentList, err := repositories.FindArticleCommentsByArticleID(db, articleID)
	if err != nil {
		return models.Article{}, err
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

func PostNiceService(article models.Article) (models.Article, error) {
	db, err := connectDB()
	if err != nil {
		return models.Article{}, err
	}
	defer db.Close()

	err = repositories.UpdateArticleNice(db, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	article, err = repositories.GetArticleDetailByID(db, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}
