# go-fasttext

Go bindings for [fastText](https://github.com/facebookresearch/fastText) language identification via a thin C wrapper and CGO.

## Requirements

- Go with CGO enabled
- C++17 compiler (`c++` / `clang++`)
- `make`, `ar`

## Build

```bash
# initialize submodule if needed
git submodule update --init --recursive

# build PIC static library used by CGO
make lib
```

This compiles `fastText/src` (excluding `main.cc`) and `cwrapper/` into `build/libgo_fasttext.a`.

## Download a LID model

Models are not vendored. Download the compressed 176-language model (~917KB):

```bash
mkdir -p models
curl -L -o models/lid.176.ftz \
  https://dl.fbaipublicfiles.com/fasttext/supervised-models/lid.176.ftz
```

Larger / alternate models: [language identification docs](https://fasttext.cc/docs/en/language-identification.html).

## Usage

```go
package main

import (
	"fmt"
	"log"

	gofasttext "go-fasttext"
)

func main() {
	if err := gofasttext.LoadModel("models/lid.176.ftz"); err != nil {
		log.Fatal(err)
	}
	defer gofasttext.Close()

	lang, conf, err := gofasttext.Predict("你好世界")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s (%.4f)\n", lang, conf)
}
```

API:

- `LoadModel(path string) error` — load once (replaces any previous model)
- `Predict(text string) (lang string, confidence float32, err error)` — top-1 language code without `__label__`
- `Close()` — free the loaded model

## Test

```bash
make lib
go test ./...
```

Tests look for `models/lid.176.ftz` (or `models/lid.176.bin`), or the path in `TEST_LID_MODEL`. If no model is present, prediction tests are skipped.
