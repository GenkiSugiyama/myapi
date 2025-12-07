package services

import (
	"database/sql"
	"errors"

	"github.com/GenkiSugiyama/myapi/apperrors"
	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/repositories"
)

func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "fail to record data")
		return models.Article{}, err
	}

	return newArticle, nil
}

func (s *MyAppService) ArticleListService(page int) ([]models.Article, error) {
	articleList, err := repositories.FindArticles(s.db, page)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")
		return nil, err
	}

	if len(articleList) == 0 {
		err := apperrors.NAData.Wrap(ErrNoData, "no data")
		return nil, err
	}

	return articleList, nil
}

func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	article, err := repositories.GetArticleDetailByID(s.db, articleID)
	if err != nil {
		// 1件もデータが取得されたなかった場合のエラーハンドリング
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NAData.Wrap(err, "no data")
			return models.Article{}, err
		}
		// それ以外はDB接続等の環境的なエラーとして扱う
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")

		return models.Article{}, err
	}

	commentList, err := repositories.FindArticleCommentsByArticleID(s.db, articleID)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")
		return models.Article{}, err
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	err := repositories.UpdateArticleNice(s.db, article.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NoTargetData.Wrap(err, "does not exist taarget article")
			return models.Article{}, err
		}
		err = apperrors.UpdateDataFailed.Wrap(err, "fail to update nice count")
		return models.Article{}, err
	}

	article, err = repositories.GetArticleDetailByID(s.db, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}
