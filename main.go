package main

import (
	"flag"
	"log"
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
		runServer()
	}
}
