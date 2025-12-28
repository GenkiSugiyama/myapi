package middlewares

import (
	"context"
	"sync"
)

var (
	logNo int = 1
	mu    sync.Mutex
)

func SetTraceID(ctx context.Context, traceID int) context.Context {
	// context.WithValueは元となるコンテキストに指定されたkey, valueのペアを追加した新しいコンテキストを返す
	return context.WithValue(ctx, "traceID", traceID)
}

func GetTraceID(ctx context.Context) int {
	// ctx.Value()で指定されたkeyに対応するvalueを取得する
	// その時のvalueの型はany型なので、元の型にアサーションする必要がある
	id := ctx.Value("traceID")

	if idInt, ok := id.(int); ok {
		return idInt
	}

	return 0
}

func newTraceID() int {
	var no int

	mu.Lock()
	no = logNo
	logNo++
	mu.Unlock()

	return no
}
