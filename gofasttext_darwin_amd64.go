//go:build darwin && amd64

package gofasttext

/*
#cgo LDFLAGS: ${SRCDIR}/libs/darwin_amd64/libgo_fasttext.a -lc++ -lm -pthread
*/
import "C"
