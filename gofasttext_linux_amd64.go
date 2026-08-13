//go:build linux && amd64

package gofasttext

/*
#cgo LDFLAGS: ${SRCDIR}/libs/linux_amd64/libgo_fasttext.a -lstdc++ -lm -pthread
*/
import "C"
