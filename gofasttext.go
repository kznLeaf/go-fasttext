package gofasttext

/*
#cgo CFLAGS: -I${SRCDIR}/cwrapper
#cgo CXXFLAGS: -std=c++17 -O3 -I${SRCDIR}/cwrapper -I${SRCDIR}/fastText/src
#cgo LDFLAGS: -L${SRCDIR}/build -lgo_fasttext -lstdc++ -lm -pthread
#include <stdlib.h>
#include "fasttext_wrapper.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

const (
	errBufLen  = 512
	langBufLen = 64
)

var mu sync.RWMutex

// LoadModel loads a fastText supervised model (e.g. lid.176.ftz) for language ID.
// Replaces any previously loaded model.
func LoadModel(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("model path is empty")
	}

	mu.Lock()
	defer mu.Unlock()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	errBuf := make([]byte, errBufLen)
	rc := C.ft_load_model(cPath, (*C.char)(unsafe.Pointer(&errBuf[0])), C.int(errBufLen))
	if rc != 0 {
		msg := C.GoString((*C.char)(unsafe.Pointer(&errBuf[0])))
		if msg == "" {
			msg = "failed to load model"
		}
		return errors.New(msg)
	}
	return nil
}

// Predict returns the top-1 language code (without "__label__") and confidence.
func Predict(text string) (lang string, confidence float32, err error) {
	if text == "" {
		return "", 0, errors.New("text is empty")
	}

	mu.RLock()
	defer mu.RUnlock()

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	langBuf := make([]byte, langBufLen)
	errBuf := make([]byte, errBufLen)
	var prob C.float

	rc := C.ft_predict(
		cText,
		(*C.char)(unsafe.Pointer(&langBuf[0])),
		C.int(langBufLen),
		&prob,
		(*C.char)(unsafe.Pointer(&errBuf[0])),
		C.int(errBufLen),
	)
	if rc != 0 {
		msg := C.GoString((*C.char)(unsafe.Pointer(&errBuf[0])))
		if msg == "" {
			msg = "prediction failed"
		}
		return "", 0, errors.New(msg)
	}

	lang = C.GoString((*C.char)(unsafe.Pointer(&langBuf[0])))
	if lang == "" {
		return "", 0, fmt.Errorf("empty language label")
	}
	return lang, float32(prob), nil
}

// Close frees the loaded model.
func Close() {
	mu.Lock()
	defer mu.Unlock()
	C.ft_free_model()
}
