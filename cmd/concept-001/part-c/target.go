package main

import (
	"fmt"
	"net/http"
)

// --- go.buildinfo ---

// go build -o target target.go
// go tool objdump -s "main.main" target

// strings target | grep "AKIA"

const awsKey = "AKIAIOSFODNN7EXAMPLE_SECRET_KEY"

func main() {
	fmt.Println("Hello, Golab!", awsKey)
	http.Get("http://example.com")
}
