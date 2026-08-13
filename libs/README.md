# Prebuilt static libraries

Vendored `libgo_fasttext.a` archives for CGO linking. Users on supported platforms can `go get` this module without running `make lib`.

## Layout

```
libs/
  linux_amd64/libgo_fasttext.a
  linux_arm64/libgo_fasttext.a
  darwin_amd64/libgo_fasttext.a
  darwin_arm64/libgo_fasttext.a
```

Directory names follow Go's `GOOS_GOARCH` convention.

## Rebuilding (maintainers)

After changing `fastText/src` or `cwrapper/`:

```bash
make lib
make lib-install   # copies build/libgo_fasttext.a to libs/$(go env GOOS)_$(go env GOARCH)/
```

CI (`.github/workflows/build-libs.yml`) rebuilds all four POSIX targets on push to `main`.

## Requirements

- C++17 compiler with `-fPIC` support
- `make` and `ar`
- fastText submodule initialized: `git submodule update --init --recursive`
