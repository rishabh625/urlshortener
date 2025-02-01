package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server to be started at 8080")
	mux := http.NewServeMux()
	mux.HandleFunc()

}
