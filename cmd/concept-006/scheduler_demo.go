package main

import (
	"fmt"
	"net/http"
	"sync"
)

// --- how to run ---
// 1. go build -o scheduler_demo scheduler_demo.go
// 2. go tool objdump -s "net/http\." ./scheduler_demo | head -n 50
// 3. go tool objdump -s "main\." ./scheduler_demo
// disable inlining with "-l"
// 4. go build -gcflags="-l" -o scheduler_demo_noinline scheduler_demo.go

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
