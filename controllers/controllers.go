package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/GenkiSugiyama/myapi/models"
	"github.com/GenkiSugiyama/myapi/services"
	"github.com/gorilla/mux"
)

type MyAppController struct {
	service *services.MyAppService
}

func NewMyAppController(s *services.MyAppService) *MyAppController {
	return &MyAppController{service: s}
}

// リクエストを受け取って任意のレスポンスを書き込むための関数型ハンドラを宣言する
func (c *MyAppController) HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func (c *MyAppController) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	// json.NewDecoder()の引数にr.Bodyを渡してBody内のJSONデータをreqArticleにデコードする
	// json.Unmarshal()の場合は、デコード対象をメモリに格納する必要があるため、バイトスライスを用意しそこに内容を格納↓
	// その後デコードする必要があった
	// デコーダを使用することで、r.Bodyというストリームから直接デコードできるようになる
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	// MyAppService構造体を取得
	// MyAppService.PostArticleService()メソッドを呼び出して記事投稿処理を実行する
	article, err := c.service.PostArticleService(reqArticle)
	if err != nil {
		http.Error(w, "failed internal exec\n", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "failed to encode json\n", http.StatusInternalServerError)
		return
	}
}

func (c *MyAppController) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
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

	articleLists, err := c.service.ArticleListService(page)
	if err != nil {
		http.Error(w, "failed internal exec\n", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(articleLists); err != nil {
		errMsg := fmt.Sprintf("fail to encode json (page: %d)", page)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (c *MyAppController) ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得するためにmux.Vars()を使用する
	// mux.Vars()はmap[string]string型を返すので、パスパラメータの名前をキーにして値を取得する
	// 取得したパスパラメータは文字列型なので、数値として扱うために変換処理を行う
	articleID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid path parameter", http.StatusBadRequest)
		return
	}

	article, err := c.service.GetArticleService(articleID)
	if err != nil {
		http.Error(w, "faild internal exec\n", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(article); err != nil {
		errMsg := fmt.Sprintf("fail to encode json (articleID: %d)", article.ID)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (c *MyAppController) PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	article, err := c.service.PostNiceService(reqArticle)
	if err != nil {
		http.Error(w, "failed internal exec\n", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "failed to encode json\n", http.StatusInternalServerError)
		return
	}
}

func (c *MyAppController) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var reqComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&reqComment); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		http.Error(w, "failed internal exec\n", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, "failed to encode json\n", http.StatusInternalServerError)
		return
	}
}
