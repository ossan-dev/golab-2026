package main

import "fmt"

// Questions:
// 1. define a prototype in Go
// 2. Go's operator precedence

// --- how to run ---
// 1. go run .

// Add declares a function implemented in assembly
// It takes two int64 parameters and returns their sum
func Add(a, b int64) int64

func main() {
	x := int64(10)
	y := int64(32)

	result := Add(x, y)

	fmt.Printf("Add(%d, %d) = %d\n", x, y, result)
}
