package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// 1. go build -o suspect suspect.go
// 2. strace -f -e trace=openat,connect,sendto ./suspect

// Try to open the /etc/passwd
// [pid 61859] openat(AT_FDCWD, "/etc/passwd", O_RDONLY|O_CLOEXEC) = 3

// Try to connect to example.com with IP: 2606:4700:10::6814:179a
// [pid 61859] connect(3, {sa_family=AF_INET6, sin6_port=htons(53), sin6_flowinfo=htonl(0), inet_pton(AF_INET6, "2606:4700:10::6814:179a", &sin6_addr), sin6_scope_id=0}, 28) = 0

func main() {
	os.ReadFile("/etc/passwd")     // sneaky file read
	http.Get("http://example.com") // sneaky network call
	time.Sleep(2 * time.Second)
	fmt.Println("Done")
}
