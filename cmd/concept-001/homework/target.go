package main

import (
	"fmt"
	"net/http"
)

// --- go.buildinfo ---

// 1. go build -ldflags="-s -w" -o target_stripped target.go
// 	(Note: -s removes the symbol table, -w removes DWARF debug info).
// 2. go version -m target_stripped
//	Yes, the `buildinfo` has survived to the strip.
// 3. go tool objdump -s "main.main" target_stripped
// No, the `pclntab` hasn't survived to the strip.
// This has deleted the Standard ELF Symbol Table, not the Go's internal ones.
// Standard ELF Symbol table is the Table of Content you can find at the beginning of the book
// The Go pcnltab is the Index you can find at the end of the book

const awsKey = "AKIAIOSFODNN7EXAMPLE_SECRET_KEY"

func main() {
	fmt.Println("Hello, Golab!", awsKey)
	http.Get("http://example.com")
}
