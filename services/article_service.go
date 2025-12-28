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
	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error

	// 別ゴルーチンでリポジトリ層の関数を動かしたい
	// 1つのチャネルで送受信できる型は一つだけ
	// リポジトリ層の関数は2つの型を返す
	// 2つの型を扱う構造体を新たに定義し、その構造体の型を元にしたチャネルを作るこで戻り値が複数ある関数をゴルーチンで扱えるようにする
	type articleResult struct {
		article models.Article
		err     error
	}
	// 自分で定義したarticleResult型のチャネルを用意する
	articleChan := make(chan articleResult)
	defer close(articleChan)

	// ゴルーチン内でリポジトリ層の関数を呼び出し、受け取った2つの戻り値を定義した構造体に詰めて
	// その構造体を元にしたチャネルに送信している
	go func(ch chan<- articleResult, db *sql.DB, articleID int) {
		article, err := repositories.GetArticleDetailByID(db, articleID)
		ch <- articleResult{article: article, err: err}
	}(articleChan, s.db, articleID)

	type commentResult struct {
		commentList *[]models.Comment
		err         error
	}
	commentChan := make(chan commentResult)
	defer close(commentChan)

	go func(ch chan<- commentResult, db *sql.DB, articleID int) {
		commentList, err := repositories.FindArticleCommentsByArticleID(db, articleID)
		ch <- commentResult{commentList: &commentList, err: err}
	}(commentChan, s.db, articleID)

	// article用のゴルーチンとcommentList用のゴルーチン、どちらが先に処理を終えるかはわからない
	// チャネルの数分select文をループして、送信の準備ができたチャネルの値から受信するようにしている
	for i := 0; i < 2; i++ {
		select {
		case ar := <-articleChan:
			article, articleGetErr = ar.article, ar.err
		case cr := <-commentChan:
			commentList, commentGetErr = *cr.commentList, cr.err
		}
	}

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
