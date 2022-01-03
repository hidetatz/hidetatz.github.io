package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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
	files := readFiles("./data/articles")
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

func readYamlFrontMatter(aa *article, line string) error {
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

func readFiles(dir string) []*os.File {
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

	return files
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

func convertArticleToHTML(a *article) string {
	return convertMarkdownToHTML(a.title, strings.Join(a.contentsMD, "\n"), &a.timestamp)
}

func convertMarkdownToHTML(title, markdown string, timestamp *time.Time) string {
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

	return generateHTMLPage(fmt.Sprintf("%s | hidetatz.io", title), contentsHTML)
}

func linkToArticle(a *article) string {
	if a.url != nil {
		// in case an url is found for the article, directly link to that url
		return a.url.String()
	}

	return fmt.Sprintf("/articles/%s/%s", a.timestamp.Format(timeformat), a.path)
}

func toHTML(md string) string {
	// return string(gfm.Markdown([]byte(markdown)))
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes)
	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank | html.FootnoteReturnLinks})
	return string(markdown.ToHTML([]byte(md), parser, renderer))

}
