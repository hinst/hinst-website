package server

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

func use(v any) {
}

func ioCloseSilently(v io.Closer) {
	ignoreError(v.Close())
}

func ioClose(v io.Closer) {
	assertError(v.Close())
}

func ifElse[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}

// Source: https://dev.to/leapcell/how-to-get-the-goroutine-id-1h5o
func getThreadId() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		panic(fmt.Errorf("cannot get goroutine id: %v", err))
	}

	return int64(id)
}
