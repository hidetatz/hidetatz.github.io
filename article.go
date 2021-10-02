package main

import (
	"bufio"
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

	"golang.org/x/net/html"
)

const timeformat = "2006/01/02"

type articleType string

const (
	blogType  articleType = "blog"
	inputType articleType = "input"
)

type lang string

const (
	ja lang = "ja"
	en lang = "en"
)

type article struct {
	typ        articleType
	title      string
	timestamp  time.Time
	fileName   string
	url        *url.URL
	lang       lang
	contentsMD []string
}

func readArticles() ([]*article, error) {
	files := readFiles("./data/articles")
	var articles []*article

	for _, file := range files {
		var aa article

		aa.fileName = filepath.Base(file.Name())

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
	case "type":
		switch val {
		case string(blogType):
			aa.typ = blogType
		case string(inputType):
			aa.typ = inputType
		default:
			return fmt.Errorf("unknown type: %s", val)
		}
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

	// in case input
	if aa.url != nil && aa.title == "" {
		aa.title = getHTMLTitle(aa.url, trimExtension(aa.fileName))
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

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func getHTMLTitle(u *url.URL, fallback string) string {
	resp, err := http.Get(u.String())
	if err != nil {
		return fallback
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fallback
	}

	title, ok := traverse(doc)
	if !ok {
		return fallback
	}

	return title
}
func trimExtension(filename string) string {
	return filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
}

// title, datetime, content
const articlePageMD = `
[<- home](/)

# %s

#### %s

%s

<a href="https://twitter.com/share?ref_src=twsrc%%5Etfw" class="twitter-share-button" data-via="dty1er1" data-show-count="false">Tweet</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

---

<div style="text-align: center;">
  <a href="/">home</a>
</div>
`

func generateArticlePageHTML(a *article) string {
	contents := strings.Join(a.contentsMD, "\n")
	if a.typ == inputType {
		contents = fmt.Sprintf("### %s\n\n%s", a.url.String(), contents)
	}

	return generateHTMLPage(fmt.Sprintf("%s | dtyler.io", a.title), fmt.Sprintf(
		articlePageMD,
		a.title,
		a.timestamp.Format(timeformat),
		contents,
	))
}
