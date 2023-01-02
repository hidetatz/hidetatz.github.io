package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/feeds"
	"github.com/snabb/sitemap"
)

const timeformat = "2006/01/02"

type lang string

const (
	ja lang = "ja"
	en lang = "en"
)

type article struct {
	title      string
	timestamp  time.Time
	fileName   string
	url        *url.URL
	lang       lang
	contentsMD []string
	path       string
}

func readArticles() ([]*article, error) {
	// readYamlFrontMatter parses the given line as yaml front matter for the article
	readYamlFrontMatter := func(aa *article, line string) error {
		splitted := strings.Split(line, ": ")
		key, val := splitted[0], splitted[1]
		switch key {
		case "timestamp":
			t, err := time.Parse("2006-01-02 15:04:05", val)
			if err != nil {
				return fmt.Errorf("cannot parse timestamp: %s", val)
			}
			aa.timestamp = t
		case "url":
			u, err := url.Parse(val)
			if err != nil {
				return fmt.Errorf("cannot parse url: %s", val)
			}
			aa.url = u
		case "lang":
			if val == string(ja) {
				aa.lang = ja
			} else {
				aa.lang = en
			}
		case "title":
			aa.title = val
		default:
			return fmt.Errorf("unknown key in yaml: %s", key)
		}

		return nil
	}

	// read articles files
	dir := "./data/articles"
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	files := []*os.File{}
	for _, info := range fileInfo {
		if info.IsDir() {
			continue
		}
		f, err := os.Open(path.Join(dir, info.Name()))
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, f)
	}

	var articles []*article

	for _, file := range files {
		var aa article

		aa.fileName = filepath.Base(file.Name())
		aa.path = filepath.Base(aa.fileName[:len(aa.fileName)-len(filepath.Ext(aa.fileName))]) // trim extension

		scanner := bufio.NewScanner(file)
		inFrontMatter := true
		for scanner.Scan() {
			line := scanner.Text()

			// assuming every article has yaml front matter
			if line == "---" {
				inFrontMatter = false
				continue
			}

			// read yaml front-matter
			if inFrontMatter {
				if err := readYamlFrontMatter(&aa, line); err != nil {
					return nil, fmt.Errorf("failed to read yaml front matter: %w in %s", err, aa.fileName)
				}
				continue
			}

			// read contents as markdown
			aa.contentsMD = append(aa.contentsMD, line)
		}

		articles = append(articles, &aa)
	}

	sort.Slice(articles, func(i, j int) bool { return articles[i].timestamp.After(articles[j].timestamp) })

	return articles, nil
}

// title, datetime, content
const articlePageMD = `
[<- home](/)

# %s

#### %s

%s
`

const twitterButton = `
<p><a href="https://twitter.com/share?ref_src=twsrc%5Etfw" class="twitter-share-button" data-via="hidetatz" data-show-count="false">Tweet</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></p>
`

func convertArticleToHTML(title, markdown string, timestamp *time.Time) string {
	ts := ""
	if timestamp != nil {
		ts = timestamp.Format(timeformat)
	}
	contentsMD := fmt.Sprintf(
		articlePageMD,
		title,
		ts,
		markdown,
	)
	contentsHTML := toHTML(contentsMD)
	// workaround: if twitter share button is embedded into articlePageMD,
	// the footnotes are placed under the button which does not look good
	contentsHTML = contentsHTML + twitterButton

	return generateHTMLPage(fmt.Sprintf("%s | hidetatz.github.io", title), contentsHTML)
}

func linkToArticle(a *article) string {
	if a.url != nil {
		// in case an url is found for the article, directly link to that url
		return a.url.String()
	}

	return fmt.Sprintf("/articles/%s/%s", a.timestamp.Format(timeformat), a.path)
}

func toHTML(md string) string {
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes)
	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank | html.FootnoteReturnLinks})
	return string(markdown.ToHTML([]byte(md), parser, renderer))

}

const indexPageMD = `
# hidetatz.github.io

## About me

I'm hidetatz, an automotive software engineer developing software-defined car at [TURING Inc.](https://www.turing-motors.com/en). Making level 5, full self-driving EV.

* [GitHub/hidetatz](https://github.com/hidetatz)
* [Twitter/hidetatz](https://twitter.com/hidetatz) (rarely appear)
* [Atom/RSS feed](/feed.xml).

---

## Projects

If you love it, give a star!

* [kubecolor](https://github.com/hidetatz/kubecolor)
  - A CLI tool which colorizes the kubectl output for readability.
  - You can read my [blog article](https://hidetatz.medium.com/colorize-kubectl-output-by-kubecolor-2c222af3163a) about it.
* [incdb](https://github.com/hidetatz/incdb)
  - Incrementally developed RDB from scratch.
* [rv](https://github.com/hidetatz/rv)
  - RISC-V software emulator.

---

## Articles

%s

Some articles are available in Japanese also.

%s

---

## Other writings

* [/distsys](/distsys.html)
  - Distributed systems learning meterials (in Japanese)

---

© 2022 Hidetatz Yaginuma. Unless otherwise noted, these posts are made available under a [Creative Commons Attribution License](https://creativecommons.org/licenses/by/4.0/).
`

func generateIndexPageHTML(articles []*article) string {
	enblogsList := ""
	jablogsList := ""
	for _, a := range articles {
		switch a.lang {
		case en:
			enblogsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, linkToArticle(a))
		case ja:
			jablogsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, linkToArticle(a))
		}
	}

	contentsHTML := toHTML(fmt.Sprintf(indexPageMD, enblogsList, jablogsList))
	return generateHTMLPage("hidetatz.github.io", contentsHTML)
}

const page = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="author" content="Hidetatz Yaginuma">
  <meta name="viewport" content="width=device-width, initial-scale=1, minimal-ui">

  <style>
    body {
      box-sizing: border-box;
      min-width: 200px;
      max-width: 980px;
      margin: 0 auto;
      padding: 45px;
    }
  </style>

  <link href="/markdown.css" rel="stylesheet"></link>
  <link href="/syntax.css" rel="stylesheet"></link>
  <script type="text/javascript" async
    src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-MML-AM_CHTML">
  </script>
</head>
<body class="markdown-body">
%s

<script src="/syntax.js"></script>
<script>hljs.highlightAll();</script>
</body>
</html>
`

func generateHTMLPage(title, content string) string {
	return fmt.Sprintf(page, title, content)
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

func genFeed(articles []*article, t time.Time, count int) string {
	if count > len(articles) {
		count = len(articles)
	}

	name := "Hidetatz Yaginuma"
	email := "hidetatz@gmail.com"
	feed := &feeds.Feed{
		Title:   fmt.Sprintf("hidetatz.github.io | %s", name),
		Link:    &feeds.Link{Href: "https://hidetatz.github.io"},
		Author:  &feeds.Author{Name: name, Email: email},
		Created: t,
		Items:   make([]*feeds.Item, count),
	}

	for i := 0; i < count; i++ {
		a := articles[i]
		feed.Items[i] = &feeds.Item{
			Title:       a.title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://hidetatz.github.io/%s", linkToArticle(a))},
			Description: "The post first appeared on hidetatz.github.io.",
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

func genSiteMap(articles []*article, t time.Time) string {
	fqdn := "hidetatz.github.io"
	sm := sitemap.New()
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s", fqdn), LastMod: &t})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/distsys.html", fqdn), LastMod: &t})

	for _, a := range articles {
		sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/%s", fqdn, linkToArticle(a)), LastMod: &a.timestamp})
	}

	buff := &bytes.Buffer{}
	sm.WriteTo(buff)
	return buff.String()
}

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

func gen() {
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
			write(convertArticleToHTML(a.title, strings.Join(a.contentsMD, "\n"), &a.timestamp), fmt.Sprintf("./docs/%s/index.html", linkToArticle(a)))
		}
	}

	// index
	write(generateIndexPageHTML(articles), "./docs/index.html")

	latestArticle := articles[0]

	// sitemap, atom feed
	write(genSiteMap(articles, latestArticle.timestamp), "./docs/sitemap.xml")
	write(genFeed(articles, latestArticle.timestamp, 20), "./docs/feed.xml")

	// 404 page
	articlesCountOn404Page := 5
	articlesFor404Page := ""
	for i := 0; i < articlesCountOn404Page; i++ {
		articlesFor404Page += fmt.Sprintf("[%s](%s)  \n", articles[i].title, linkToArticle(articles[i]))
	}
	write(convertArticleToHTML("404", fmt.Sprintf(notFoundPage, articlesFor404Page), nil), "./docs/404.html")

	// other writings
	write(convertArticleToHTML("Learn distributed systems", distsys, nil), "./docs/distsys.html")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	serve := flag.Bool("serve", false, "generate and serve site on local")
	flag.Parse()

	removeAllFiles("./docs/")
	gen()
	if *serve {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		go func() {
			for {
				select {
				case <-watcher.Events:
					log.Println("regenerating...")
					removeAllFiles("./docs/")
					gen()
				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()

		err = watcher.Add("./data")
		if err != nil {
			log.Fatal(err)
		}

		err = watcher.Add("./data/articles")
		if err != nil {
			log.Fatal(err)
		}

		server := &http.Server{Addr: ":8090", Handler: http.FileServer(http.Dir("./docs"))}
		log.Printf("Serving at localhost:8090\n")
		server.ListenAndServe()
	}
}
