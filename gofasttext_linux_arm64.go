//go:build linux && arm64

package gofasttext

/*
#cgo LDFLAGS: ${SRCDIR}/libs/linux_arm64/libgo_fasttext.a -lstdc++ -lm -pthread
*/
import "C"
