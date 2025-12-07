package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/GenkiSugiyama/myapi/apperrors"
	"github.com/GenkiSugiyama/myapi/controllers/services"
	"github.com/GenkiSugiyama/myapi/models"
	"github.com/gorilla/mux"
)

type ArticleController struct {
	service services.ArticleServicer
}

func NewArticleController(s services.ArticleServicer) *ArticleController {
	return &ArticleController{service: s}
}

func (c *ArticleController) HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func (c *ArticleController) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	// json.NewDecoder()の引数にr.Bodyを渡してBody内のJSONデータをreqArticleにデコードする
	// json.Unmarshal()の場合は、デコード対象をメモリに格納する必要があるため、バイトスライスを用意しそこに内容を格納↓
	// その後デコードする必要があった
	// デコーダを使用することで、r.Bodyというストリームから直接デコードできるようになる
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	// MyAppService構造体を取得
	// MyAppService.PostArticleService()メソッドを呼び出して記事投稿処理を実行する
	article, err := c.service.PostArticleService(reqArticle)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	if err = json.NewEncoder(w).Encode(article); err != nil {
		err = apperrors.StructEncodeFailed.Wrap(err, "bad struct")
		apperrors.ErrorHandler(w, r, err)
		return
	}
}

func (c *ArticleController) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	// *URL.Query()はクエリパラーメータのKeyとKeyに対応するValueを持つmap[string]][]string型を返す
	queryMap := r.URL.Query()

	var page int
	// クエリパラメーターのキーに対応する文字列型のスライスをpに格納する、取得できたら第二変数にtrueが返ってくる
	if p, ok := queryMap["page"]; ok && len(p) > 0 {
		var err error
		page, err = strconv.Atoi(p[0])
		if err != nil {
			err = apperrors.BadParam.Wrap(err, "queryparam must be number")
			apperrors.ErrorHandler(w, r, err)
			return
		}
	} else {
		page = 1
	}

	articleLists, err := c.service.ArticleListService(page)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(articleLists); err != nil {
		err = apperrors.StructEncodeFailed.Wrap(err, "bad struct")
		apperrors.ErrorHandler(w, r, err)
		return
	}
}

func (c *ArticleController) ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得するためにmux.Vars()を使用する
	// mux.Vars()はmap[string]string型を返すので、パスパラメータの名前をキーにして値を取得する
	// 取得したパスパラメータは文字列型なので、数値として扱うために変換処理を行う
	articleID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		err = apperrors.BadParam.Wrap(err, "queryparam must be number")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	article, err := c.service.GetArticleService(articleID)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(article); err != nil {
		err = apperrors.StructEncodeFailed.Wrap(err, "bad struct")
		apperrors.ErrorHandler(w, r, err)
		return
	}
}

func (c *ArticleController) PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	article, err := c.service.PostNiceService(reqArticle)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	if err = json.NewEncoder(w).Encode(article); err != nil {
		err = apperrors.StructEncodeFailed.Wrap(err, "bad struct")
		apperrors.ErrorHandler(w, r, err)
		return
	}
}
