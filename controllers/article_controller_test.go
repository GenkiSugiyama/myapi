package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestArticleListHandler(t *testing.T) {
	var tests = []struct {
		name       string
		query      string
		resultCode int
	}{
		{name: "number query", query: "1", resultCode: http.StatusOK},
		{name: "alphabet query", query: "aaa", resultCode: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// handler用の引数の作成
			// httptest.NewRequest()でhttp.Requestのモックを作成する
			// 第一引数にHTTPメソッド、第二引数にテスト対象となるURL、第三引数にリクエストボディを指定する
			url := fmt.Sprintf("http://localhost:8080/article/list?page=%s", tt.query)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			// httptest.NewRecorder()でhttp.ResponseWriterのモックを作成する
			// ResponseRecorderをハンドラに渡すことで、ハンドラが書き込んだレスポンスをテストコード内で取得できるようになる
			res := httptest.NewRecorder()

			aCon.ArticleListHandler(res, req)

			if res.Code != tt.resultCode {
				t.Errorf("unexpected StatusCode: want %d but %d\n", tt.resultCode, res.Code)
			}
		})
	}
}

func TestArticleDetailHandler(t *testing.T) {
	var tests = []struct {
		name       string
		articleID  string
		resultCode int
	}{
		{name: "number pathparam", articleID: "1", resultCode: http.StatusOK},
		{name: "alphabet pathparam", articleID: "aaa", resultCode: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/article/%s", tt.articleID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()

			// ArticleDetailHandlerはmux.Vars()を使用してパスパラメータを取得しているため、
			// テストコード内でパスパラメータを正しく取得できるようにするためにmux.Routerを使用してハンドラを呼び出す必要がある
			// 実際のハンドラ登録と同様にmux.Routerを作成し、パスパラメータ付きのURLパターンでハンドラを登録する
			r := mux.NewRouter()
			r.HandleFunc("/article/{id:[0-9]+}", aCon.ArticleDetailHandler).Methods(http.MethodGet)
			r.ServeHTTP(res, req)

			if res.Code != tt.resultCode {
				t.Errorf("unexpected StatusCode: want %d but %d\n", tt.resultCode, res.Code)
			}
		})
	}
}
