package main

import (
	"fmt"
	"net/http"
)

// --- go.buildinfo ---

// go build -o target target.go
// go version -m target

// target: go1.25.5
// path    command-line-arguments
// build   -buildmode=exe
// build   -compiler=gc
// build   CGO_ENABLED=1
// build   CGO_CFLAGS=
// build   CGO_CPPFLAGS=
// build   CGO_CXXFLAGS=
// build   CGO_LDFLAGS=
// build   GOARCH=amd64
// build   GOOS=linux
// build   GOAMD64=v

func main() {
	fmt.Println("Hello, Golab!")
	http.Get("http://example.com")
}
