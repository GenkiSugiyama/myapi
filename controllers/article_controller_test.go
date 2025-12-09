package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
