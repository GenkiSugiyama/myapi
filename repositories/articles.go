package repositories

import (
	"database/sql"
	"fmt"

	"github.com/GenkiSugiyama/myapi/models"
)

const ArticlesPerPage = 5

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

func FindArticles(db *sql.DB, page int) ([]models.Article, error) {
	// 先にクエリ文字列を定義する
	// 値を動的に設定したい箇所は?でプレースホルダを設定する（MySQL）
	const sqlStr = `
		SELECT title, contents, username, nice
		FROM articles
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?;
	`

	// db.Queryでクエリを実行し、結果をrowsに格納する
	// 第二引数以降でプレースホルダに値を設定する
	rows, err := db.Query(sqlStr, ArticlesPerPage, (page-1)*ArticlesPerPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]models.Article, ArticlesPerPage)
	// rows.Next()で読み出すレコードがなくなるまで構造体への格納処理を繰り返す
	for rows.Next() {
		// 格納先となる構造体を用意し、rows.Scan()でフィールドに値をセットする
		// この時ポインタ型のフィールドを渡しているのは、用意したメモリ領域に値を書き込むため
		var article models.Article
		rows.Scan(&article.Title, &article.Contents, &article.UserName, &article.NiceNum)
		articles = append(articles, article)
	}
	return articles, nil
}

func GetArticleDetailByID(db *sql.DB, articleID int) (models.Article, error) {
	const sqlStr = `
		SELECT title, contents, username, nice, created_at
		FROM articles
		WHERE id = ?
	`

	// クエリ結果が0 or 1件の場合はQueryRowを使用する
	row := db.QueryRow(sqlStr, articleID)
	// 戻り値のRow型はクエリ結果が0件の場合にエラーを返すため、Errメソッドで確認する
	if err := row.Err(); err != nil {
		return models.Article{}, err
	}

	var article models.Article
	var createdAt sql.NullTime
	// rows.Scan, row.ScanではNull許可のカラムのNull値をそのままフィールドにセットできない
	//　Nullかもしれない値を受け取るために、sql.NullTimeでcreated_atの値を受け取る
	// sql.NullTime.ValidでNullかどうかを判定し、Validがtrueの場合にのみTimeフィールドの値を構造体のフィールドにセットする
	err := row.Scan(&article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdAt)
	if err != nil {
		return models.Article{}, err
	}

	if createdAt.Valid {
		article.CreatedAt = createdAt.Time
	}

	return article, nil
}

func UpdateArticleNice(db *sql.DB, articleID int) error {
	// 「SELECTで更新前のniceの値を取得し、UPDATEで更新する」という1つの処理で2回のクエリを
	// 実行するため、db.Begin()でトランザクションをはる
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	const selectQueryStr = `
		SELECT nice
		FROM articles
		WHERE id = ?;
	`

	// クエリの実施やスキャン処理でエラーが発生した場合は、tx.Rollback()でトランザクションを取り消す
	row := tx.QueryRow(selectQueryStr, articleID)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}

	var niceNum int
	err = row.Scan(&niceNum)
	if err != nil {
		tx.Rollback()
		return err
	}

	const updateQueryStr = `
		Update articles
		SET nice = ?
		WHERE id = ?;
	`
	_, err = tx.Exec(updateQueryStr, niceNum+1, articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
