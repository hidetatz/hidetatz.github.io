package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	cname = "hidetatz.io"

	// the number of "recent articles" on 404 page
	articlesCountOn404Page = 5
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

	//go:embed data/404.md
	notFoundPage string

	//go:embed data/robots.txt
	robotsTxt string
)

func gen() {
	// required for GitHub pages
	write(cname, "./docs/CNAME")

	// robots.txt
	write(robotsTxt, "./docs/robots.txt")

	// assets
	write(markdownCSS, "./docs/markdown.css")
	write(syntaxCSS, "./docs/syntax.css")
	write(syntaxJS, "./docs/syntax.js")
	write(favicon, "./docs/favicon.ico")

	// read articles
	articles, err := readArticles()
	if err != nil {
		log.Fatal(err)
	}

	// write each articles
	for _, a := range articles {
		if a.url == nil || a.typ == inputType {
			write(generateArticlePageHTML(a), fmt.Sprintf("./docs/%s/index.html", link(a)))
		}
	}

	// index
	write(generateIndexPageHTML(articles), "./docs/index.html")

	latestArticle := articles[0]

	// sitemap, atom feed
	write(genSiteMap(articles, latestArticle.timestamp, cname), "./docs/sitemap.xml")
	write(genAtom(articles, latestArticle.timestamp, 20, cname), "./docs/feed.xml")

	// 404 page
	articlesFor404Page := ""
	for i := 0; i < articlesCountOn404Page; i++ {
		articlesFor404Page += fmt.Sprintf("[%s](%s)  \n", articles[i].title, link(articles[i]))
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
