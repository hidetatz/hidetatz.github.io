package main

import (
	"flag"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	g := flag.Bool("gen", false, "generate static site from markdown articles")
	n := flag.Bool("new", false, "generate a new article file template")

	flag.Parse()

	switch {
	case *g:
		removeAllFiles("./docs/")
		gen()
	case *n:
		newFile()
	default:
		removeAllFiles("./docs/")
		gen()
		runServer()
	}
}
