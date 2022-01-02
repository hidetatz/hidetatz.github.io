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

	gfm "github.com/shurcooL/github_flavored_markdown"
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
			continue // exclude ja/
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

const articleFooter = `
<p><a href="https://twitter.com/share?ref_src=twsrc%5Etfw" class="twitter-share-button" data-via="hidetatz" data-show-count="false">Tweet</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></p>
`

func convertArticleToHTML(a *article) string {
	contentsMD := fmt.Sprintf(
		articlePageMD,
		a.title,
		a.timestamp.Format(timeformat),
		strings.Join(a.contentsMD, "\n"),
	)
	contentsHTML := toHTML(contentsMD)

	return generateHTMLPage(fmt.Sprintf("%s | hidetatz.io", a.title), contentsHTML) + articleFooter
}

func convertMarkdownToHTML(title, markdown string) string {
	contentsMD := fmt.Sprintf(
		articlePageMD,
		title,
		"", // pass empty to not show timestamp. This looks too lazy but this is sufficient for now.
		markdown,
	)
	contentsHTML := toHTML(contentsMD)

	return generateHTMLPage(fmt.Sprintf("%s | hidetatz.io", title), contentsHTML) + articleFooter
}

func linkToArticle(a *article) string {
	if a.url != nil {
		// in case an url is found for the article, directly link to that url
		return a.url.String()
	}

	return fmt.Sprintf("/articles/%s/%s", a.timestamp.Format(timeformat), a.path)
}

func toHTML(markdown string) string {
	return string(gfm.Markdown([]byte(markdown)))
}
