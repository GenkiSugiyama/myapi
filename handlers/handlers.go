package handlers

import (
	"encoding/json"
	"errors"
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
	// リクエストボディを受け取るためのバイトスライスを用意する
	// ヘッダーのContent-lengthからリクエストボディの長さを取得し、
	// その長さ分のバイトスライスを作成する
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		http.Error(w, "fail to get content length\n", http.StatusBadRequest)
		return
	}
	reqBodybuffer := make([]byte, length)

	// Request.Body.Read()でリクエストボディの内容をreqBodybufferに読み込む
	// Read()はファイルの読み込みが完了した際にio.EOFエラーを返すため、errors.Is()でio.EOFかどうかを確認し
	// io.EOF以外のエラーの場合は失敗とみなして500エラーを返す
	if _, err := r.Body.Read(reqBodybuffer); !errors.Is(err, io.EOF) {
		http.Error(w, "fail to get request body\n", http.StatusInternalServerError)
		return
	}
	// Read()で読み込んだボディはCloseする必要があるので必ず閉じるためにdeferを使ってメソッドを呼び出している
	defer r.Body.Close()

	var reqArticle models.Article
	if err := json.Unmarshal(reqBodybuffer, &reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	article := reqArticle
	jsonData, err := json.Marshal(article)
	if err != nil {
		http.Error(w, "fail to encode json\n", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
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
	jsonData, err := json.Marshal(articleLists)
	if err != nil {
		errMsg := fmt.Sprintf("fail to encode json (page: %d)", page)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
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

	articleList := []models.Article{models.Article1, models.Article2}
	var targetArticle models.Article
	for _, article := range articleList {
		if article.ID == articleID {
			targetArticle = article
			break
		}
	}

	jsonData, err := json.Marshal(targetArticle)
	if err != nil {
		errMsg := fmt.Sprintf("fail to encode json (articleID: %d)", articleID)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Posting Nice...\n")
	jsonData, err := json.Marshal(models.Article1)
	if err != nil {
		http.Error(w, "fail to encode json", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(models.Comment1)
	if err != nil {
		http.Error(w, "fail to encode json", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
