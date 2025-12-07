package services

import (
	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/repositories"
)

func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		return models.Article{}, err
	}

	return newArticle, nil
}

func (s *MyAppService) ArticleListService(page int) ([]models.Article, error) {
	articleList, err := repositories.FindArticles(s.db, page)
	if err != nil {
		return nil, err
	}

	return articleList, nil
}

func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	article, err := repositories.GetArticleDetailByID(s.db, articleID)
	if err != nil {
		return models.Article{}, err
	}

	commentList, err := repositories.FindArticleCommentsByArticleID(s.db, articleID)
	if err != nil {
		return models.Article{}, err
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	err := repositories.UpdateArticleNice(s.db, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	article, err = repositories.GetArticleDetailByID(s.db, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}
