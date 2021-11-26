package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	g := flag.Bool("gen", false, "generate static site from markdown articles")

	flag.Parse()

	switch {
	case *g:
		removeAllFiles("./docs/")
		gen()
	default:
		removeAllFiles("./docs/")
		gen()
		server := &http.Server{Addr: ":8080", Handler: http.FileServer(http.Dir("./docs"))}
		fmt.Printf("Serving at localhost:8080\n")
		server.ListenAndServe()
	}
}
