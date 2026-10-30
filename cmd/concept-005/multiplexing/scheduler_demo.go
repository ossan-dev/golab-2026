package main

import (
	"fmt"
	"net/http"
	"sync"
)

// --- how to run ---
// 1. go build -o scheduler_demo scheduler_demo.go
// 2. strace -f -e trace=connect,clone,clone3 ./scheduler_demo

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d starting...\n", id)

			// force a syscall (network connection)
			http.Get("http://example.com")

			fmt.Printf("Goroutine %d finished.\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All done!")
}
