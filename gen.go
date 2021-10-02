package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	cname = "dtyler.io"

	// the number of "recent articles" on 404 page
	articlesCountOn404Page = 7
)

var (
	//go:embed assets/favicon.ico
	favicon string

	//go:embed assets/markdown.css
	markdownCSS string

	//go:embed assets/syntax.css
	syntaxCSS string

	//go:embed assets/highlight.pack.js
	syntaxJS string

	//go:embed data/about.md
	about string

	//go:embed data/404.md
	notFoundPage string

	//go:embed data/robots.txt
	robotsTxt string
)

func gen() {
	// required for GitHub pages
	write(cname, "./docs/CNAME")
	write(robotsTxt, "./docs/robots.txt")

	// assets
	write(markdownCSS, "./docs/markdown.css")
	write(syntaxCSS, "./docs/syntax.css")
	write(syntaxJS, "./docs/syntax.js")
	write(favicon, "./docs/favicon.ico")

	articles := ReadArticles("./data/articles")
	articleList := ListArticlesHref(articles)

	for _, a := range articles {
		// if url != nil, no need to generate the page because it is linked from nowhere
		if a.URL == nil {
			write(GenerateArticlePageHTML(a, true), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.FormatTime(), a.FileNameWithoutExtension()))
		}
	}

	articlesJA := ReadArticles("./data/articles/ja")
	articlesJAList := ListArticlesHref(articlesJA)

	for _, a := range articlesJA {
		if a.URL == nil {
			write(GenerateArticlePageHTML(a, false), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.FormatTime(), a.FileNameWithoutExtension()))
		}
	}

	idx := generateIndexPageHTML(strings.Join(articleList, "\n"))
	idxJA := generateJaIndexPageHTML(strings.Join(articlesJAList, "\n"))

	write(idx, "./docs/index.html")
	write(idxJA, "./docs/ja/index.html")

	write(genSiteMap(append(articles, articlesJA...), cname), "./docs/sitemap.xml")
	write(genAtom(articles, 20, cname), "./docs/feed.xml")
	write(genAtom(articlesJA, 20, cname), "./docs/feed_ja.xml")

	write(generateHTMLPage("about | dtyler.io", about), "./docs/about/index.html")

	inputs := readInputs("./data/inputs")
	inputList := listInputsHref(inputs)
	inputidx := generateIndexPageHTML(strings.Join(inputList, "\n"))
	for _, i := range inputs {
		write(generateInputPageHTML(i), fmt.Sprintf("./docs/inputs/%s/%s/index.html", i.FormatTime(), i.FileNameWithoutExtension()))
	}
	write(inputidx, "./docs/inputs/index.html")

	articlesFor404Page := ""
	for i := 0; i < articlesCountOn404Page; i++ {
		articlesFor404Page += articleList[i] + "\n"
	}
	write(generateHTMLPage("404 | dtyler.io", fmt.Sprintf(notFoundPage, articlesFor404Page)), "./docs/404.html")
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
