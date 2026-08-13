//go:build darwin && arm64

package gofasttext

/*
#cgo LDFLAGS: ${SRCDIR}/libs/darwin_arm64/libgo_fasttext.a -lc++ -lm -pthread
*/
import "C"
