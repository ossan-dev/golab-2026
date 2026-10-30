package main

import (
	"fmt"
	"net/http"
)

// --- go.buildinfo ---

// go build -o target target.go
// go tool objdump -s "main.main" target

func main() {
	fmt.Println("Hello, Golab!")
	http.Get("http://example.com")
}
