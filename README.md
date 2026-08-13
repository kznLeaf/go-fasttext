# go-fasttext

> [!NOTE]
> This project only binds the fastText `predict` / `predictLine` APIs.
> Windows platform is NOT supported.

Go bindings for [fastText](https://github.com/facebookresearch/fastText) language identification via a thin C wrapper and CGO.

Prebuilt static libraries for **linux/darwin (amd64, arm64)** are vendored under `libs/`, so you can use the package after `go get` without compiling fastText.

Darwin prebuilds target **macOS 12+** (`MACOSX_DEPLOYMENT_TARGET=12.0`) so they link cleanly on older hosts.

## Requirements

- Go with CGO enabled (`CGO_ENABLED=1`, the default)
- C/C++ linker (`gcc`/`clang` with C++ standard library) for the final link step

Supported platforms: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.

On other platforms, build from source (see [Building from source](#building-from-source)).

## Install

```bash
go get github.com/kznLeaf/go-fasttext
```

No `make lib` step is required on supported platforms.

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

	gofasttext "github.com/kznLeaf/go-fasttext"
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
go test ./...
```

Tests look for `models/lid.176.ftz` (or `models/lid.176.bin`), or the path in `TEST_LID_MODEL`. If no model is present, prediction tests are skipped.

## Building from source

Required when modifying `fastText/src` or `cwrapper/`, or on unsupported platforms:

```bash
git submodule update --init --recursive
make lib
make lib-install   # copies to ./libs/$(go env GOOS)_$(go env GOARCH)/
```

CI rebuilds all four platform libraries on changes to native code (see `.github/workflows/build-libs.yml`).

## Limitations

- CGO must stay enabled; `CGO_ENABLED=0` is not supported.
- Static libraries are platform-specific; cross-compiling requires a matching prebuilt `libs/{GOOS}_{GOARCH}/libgo_fasttext.a`.
- Windows is not included in vendored builds.
