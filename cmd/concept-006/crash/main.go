package main

import "fmt"

func main() {
	a := 10
	b := 0
	c := divide(a, b)
	fmt.Println(c)
}

func divide(x, y int) int {
	return x / y
}
