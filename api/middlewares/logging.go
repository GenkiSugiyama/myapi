package middlewares

import (
	"log"
	"net/http"
)

type resLoggingWriter struct {
	// interfaceを埋め込むと埋め込まれた構造体はそのinterfaceを実装したことになる
	// http.ResponseWriterを受け取るメソッドや関数に対してresLoggingWriter構造体を渡せるようになる
	http.ResponseWriter
	code int
}

func NewResLoggingWriter(w http.ResponseWriter) *resLoggingWriter {
	return &resLoggingWriter{ResponseWriter: w, code: http.StatusOK}
}

// interfaceメソッドのオーバーライド
// http.ResponseWriterを埋め込んだ状態でinterfaceで定義済みのメソッドを構造体でも定義すると、埋め込んだinterfaceのメソッドがオーバーライドされる
// resLogginWriter.WriteHeader()ではこの処理が呼び出され、resLoggintWriter.ResponseWriter.WriteHeader()で
// New関数で埋め込まれたhttp.ResponseWriter.WriteHeader()が呼び出される
// オーバーライドされないメソッドは、resLoggingWriter.Write()とresLoggingWriter.ResponseWriter.Write()が同じメソッドとして扱われる
func (rlw *resLoggingWriter) WriteHeader(code int) {
	rlw.code = code
	rlw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		traceID := newTraceID()

		log.Printf("[%d]%s %s\n", traceID, req.RequestURI, req.Method)

		// req.Context()でリクエストのコンテキストを取得し、そこにtraceIDをセットした新しいコンテキストを作成
		ctx := SetTraceID(req.Context(), traceID)
		// req.WithContext()で元のリクエストに新しいコンテキストをセットした新しいリクエストを作成
		req = req.WithContext(ctx)
		rlw := NewResLoggingWriter(w)

		next.ServeHTTP(rlw, req)

		log.Printf("[%d] res: %d\n", traceID, rlw.code)
	})
}
