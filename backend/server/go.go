package server

import "io"

func use(v any) {
}

func ioCloseSilently(v io.Closer) {
	ignoreError(v.Close())
}

func ioClose(v io.Closer) {
	assertError(v.Close())
}
