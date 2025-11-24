package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/GenkiSugiyama/myapi/models"
	"github.com/gorilla/mux"
)

// リクエストを受け取って任意のレスポンスを書き込むための関数型ハンドラを宣言する
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	// json.NewDecoder()の引数にr.Bodyを渡してBody内のJSONデータをreqArticleにデコードする
	// json.Unmarshal()の場合は、デコード対象をメモリに格納する必要があるため、バイトスライスを用意しそこに内容を格納↓
	// その後デコードする必要があった
	// デコーダを使用することで、r.Bodyというストリームから直接デコードできるようになる
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	article := reqArticle

	json.NewEncoder(w).Encode(article)
}

func ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	// *URL.Query()はクエリパラーメータのKeyとKeyに対応するValueを持つmap[string]][]string型を返す
	queryMap := r.URL.Query()

	var page int
	// クエリパラメーターのキーに対応する文字列型のスライスをpに格納する、取得できたら第二変数にtrueが返ってくる
	if p, ok := queryMap["page"]; ok && len(p) > 0 {
		var err error
		page, err = strconv.Atoi(p[0])
		if err != nil {
			http.Error(w, "Invalide query parameter", http.StatusBadRequest)
			return
		}
	} else {
		page = 1
	}

	articleLists := []models.Article{models.Article1, models.Article2}

	if err := json.NewEncoder(w).Encode(articleLists); err != nil {
		errMsg := fmt.Sprintf("fail to encode json (page: %d)", page)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得するためにmux.Vars()を使用する
	// mux.Vars()はmap[string]string型を返すので、パスパラメータの名前をキーにして値を取得する
	// 取得したパスパラメータは文字列型なので、数値として扱うために変換処理を行う
	articleID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid path parameter", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(models.Article1); err != nil {
		errMsg := fmt.Sprintf("fail to encode json (articleID: %d)", articleID)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// articleList := []models.Article{models.Article1, models.Article2}
	// var targetArticle models.Article
	// for _, article := range articleList {
	// 	if article.ID == articleID {
	// 		targetArticle = article
	// 		break
	// 	}
	// }

	// if err := json.NewEncoder(w).Encode(targetArticle); err != nil {
	// 	errMsg := fmt.Sprintf("fail to encode json (articleID: %d)", articleID)
	// 	http.Error(w, errMsg, http.StatusInternalServerError)
	// 	return
	// }
}

func PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(reqArticle)
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var reqComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&reqComment); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(reqComment)
}
