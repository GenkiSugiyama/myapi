package services

import (
	"database/sql"
	"errors"
	"sync"

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
	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error

	var aMux sync.Mutex
	var cMux sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2)

	// メインゴルーチンで定義した変数を直接別のゴルーチンで参照してしまうのは競合状態になる可能性があるため避けるべき
	go func(db *sql.DB, articleID int) {
		defer wg.Done()
		aMux.Lock()
		article, articleGetErr = repositories.GetArticleDetailByID(db, articleID)
		aMux.Unlock()
	}(s.db, articleID)

	go func(db *sql.DB, articleID int) {
		defer wg.Done()
		cMux.Lock()
		commentList, commentGetErr = repositories.FindArticleCommentsByArticleID(db, articleID)
		cMux.Unlock()
	}(s.db, articleID)

	wg.Wait()

	if articleGetErr != nil {
		// 1件もデータが取得されたなかった場合のエラーハンドリング
		if errors.Is(articleGetErr, sql.ErrNoRows) {
			err := apperrors.NAData.Wrap(articleGetErr, "no data")
			return models.Article{}, err
		}
		// それ以外はDB接続等の環境的なエラーとして扱う
		err := apperrors.GetDataFailed.Wrap(articleGetErr, "fail to get data")
		return models.Article{}, err
	}

	if commentGetErr != nil {
		err := apperrors.GetDataFailed.Wrap(commentGetErr, "fail to get data")
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
