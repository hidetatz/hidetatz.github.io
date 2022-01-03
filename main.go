package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/snabb/sitemap"
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

	//go:embed data/distsys.md
	distsys string

	//go:embed data/inputs.md
	inputs string

	//go:embed data/robots.txt
	robotsTxt string
)

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

func genAtom(articles []*article, t time.Time, count int, fqdn string) string {
	if count > len(articles) {
		count = len(articles)
	}

	name := "Hidetatz Yaginuma"
	email := "hidetatz@gmail.com"
	feed := &feeds.Feed{
		Title:   fmt.Sprintf("hidetatz.io | %s", name),
		Link:    &feeds.Link{Href: "https://hidetatz.io"},
		Author:  &feeds.Author{Name: name, Email: email},
		Created: t,
		Items:   make([]*feeds.Item, count),
	}

	for i := 0; i < count; i++ {
		a := articles[i]
		feed.Items[i] = &feeds.Item{
			Title:       a.title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s/%s", fqdn, linkToArticle(a))},
			Description: "The post first appeared on hidetatz.io.",
			Author:      &feeds.Author{Name: name, Email: email},
			Created:     a.timestamp,
		}

	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	return atom
}

func genSiteMap(articles []*article, t time.Time, fqdn string) string {
	sm := sitemap.New()
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s", fqdn), LastMod: &t})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/ja/", fqdn), LastMod: &t})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/about/", fqdn), LastMod: &t})

	for _, a := range articles {
		sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/%s", fqdn, linkToArticle(a)), LastMod: &a.timestamp})
	}

	buff := &bytes.Buffer{}
	sm.WriteTo(buff)
	return buff.String()
}

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
		if a.url == nil {
			write(convertArticleToHTML(a), fmt.Sprintf("./docs/%s/index.html", linkToArticle(a)))
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
		articlesFor404Page += fmt.Sprintf("[%s](%s)  \n", articles[i].title, linkToArticle(articles[i]))
	}
	write(generateHTMLPage("404 | hidetatz.io", fmt.Sprintf(notFoundPage, articlesFor404Page)), "./docs/404.html")

	// writings
	write(convertMarkdownToHTML("Learn distributed systems", distsys, nil), "./docs/distsys.html")
}

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
