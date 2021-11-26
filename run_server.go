package main

import (
	"fmt"
	"net/http"
)

func runServer() {
	server := &http.Server{Addr: ":8080", Handler: http.FileServer(http.Dir("./docs"))}

	fmt.Printf("Serving at localhost:8080\n")
	server.ListenAndServe()
}
