package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	cname = "dtyler.io"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	removeAllFiles("./docs/")
	gen()
}

func gen() {
	write(cname, "./docs/CNAME") // required for GitHub pages
	write(css, "./docs/markdown.css")
	cp("favicon.ico", "./docs/favicon.ico")
	write(GenerateHTMLPage("about", About), "./docs/about/index.html")

	articles := ReadArticles("./data/articles")
	articleList := ListArticlesHref(articles)

	for _, a := range articles {
		// if url != nil, no need to generate the page because it is linked from nowhere
		if a.URL == nil {
			write(GenerateArticlePageHTML(a), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.FormatTime(), a.FileNameWithoutExtension()))
		}
	}

	articlesJA := ReadArticles("./data/articles/ja")
	articlesJAList := ListArticlesHref(articlesJA)

	for _, a := range articlesJA {
		if a.URL == nil {
			write(GenerateArticlePageHTML(a), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.FormatTime(), a.FileNameWithoutExtension()))
		}
	}

	idx := GenerateIndexPageHTML(strings.Join(articleList, "\n"))
	idxJA := GenerateJaIndexPageHTML(strings.Join(articlesJAList, "\n"))

	write(idx, "./docs/index.html")
	write(idxJA, "./docs/ja/index.html")

}

func cp(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal(err)
	}
}

func write(content, fileNameWithDir string) {
	err := os.MkdirAll(filepath.Dir(fileNameWithDir), 0775)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(fileNameWithDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

func removeAllFiles(dir string) {
	d, err := os.Open(dir)
	defer func() {
		_ = d.Close()
	}()

	if err != nil {
		log.Fatal(err)
	}

	fileinfo, err := d.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileinfo {
		if info.IsDir() {
			removeAllFiles(filepath.Join(dir, info.Name()))
		}

		err := os.Remove(filepath.Join(dir, info.Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
