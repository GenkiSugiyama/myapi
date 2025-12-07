package services

import "database/sql"

// コントローラ構造体を初期化するために必要なサービス構造体の実体を定義する
// article_service.goでArticleServicerインターフェース、comment_service.goでCommentServicerインターフェースを実装している
type MyAppService struct {
	db *sql.DB
}

func NewMyAppService(db *sql.DB) *MyAppService {
	return &MyAppService{db: db}
}
