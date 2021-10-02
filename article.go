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
)

const timeformat = "2006/01/02"

// title, datetime, content
const articlePageMD = `
[<- home](%s)

# %s

#### %s

%s

<a href="https://twitter.com/share?ref_src=twsrc%%5Etfw" class="twitter-share-button" data-via="dty1er1" data-show-count="false">Tweet</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

---

<div style="text-align: center;">
  <a href="%s">home</a>
</div>
`

type Article struct {
	FileName   string
	Title      string
	Timestamp  time.Time
	URL        *url.URL // non-nil if the article is extenal link
	ContentsMD []string
}

// ReadArticles reads articles on the given directory
// then return all the articles ordered by the timestamp desc
func ReadArticles(dir string) []*Article {
	files := readFiles(dir)

	articles := []*Article{}
	for _, file := range files {
		scanner := bufio.NewScanner(file)
		firstLine := true
		lines := []string{}
		article := &Article{}
		for scanner.Scan() {
			if firstLine {
				if err := decodeFirstLine(article, scanner.Text()); err != nil {
					log.Fatalf("failed to decode the article first line: %s, %s", file.Name(), err)
				}
				firstLine = false
				continue
			}
			lines = append(lines, scanner.Text())
		}
		article.FileName = filepath.Base(file.Name())
		article.ContentsMD = lines
		articles = append(articles, article)
	}

	sort.Slice(articles, func(i, j int) bool { return articles[i].Timestamp.After(articles[j].Timestamp) })
	return articles
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

// parseFirstLine parses the first line of an article file.
// The format must be
// title of article---2019-04-01 12:30:00
// or
// title of article---2019-04-01 12:30:00---https://url.to.article
func decodeFirstLine(article *Article, firstLine string) error {
	splitted := strings.Split(firstLine, "---")
	if len(splitted) != 2 && len(splitted) != 3 {
		return fmt.Errorf("invalid format of article first line: %s", firstLine)
	}

	article.Title = splitted[0]

	t, err := time.Parse("2006-01-02 15:04:05", splitted[1])
	if err != nil {
		return fmt.Errorf("Invalid time format: '%s'", splitted[1])
	}

	article.Timestamp = t

	// if no url found, set only title and timestamp then return
	if len(splitted) != 3 {
		return nil
	}

	u, err := url.Parse(splitted[2])
	if err != nil {
		return fmt.Errorf("failed to parse URL: '%s'", splitted[2])
	}

	article.URL = u
	return nil
}

func ListArticlesHref(articles []*Article) []string {
	ret := make([]string, len(articles))
	for i, article := range articles {
		ret[i] = article.Link()
	}

	return ret
}

func (a *Article) FileNameWithoutExtension() string {
	return filepath.Base(a.FileName[:len(a.FileName)-len(filepath.Ext(a.FileName))])
}

func (a *Article) Link() string {
	formatedTime := a.FormatTime()
	if a.URL != nil {
		// the article is external URL
		return fmt.Sprintf(
			`<a href="%s">%s - %s</a><br>`,
			a.URL.String(),
			formatedTime,
			a.Title,
		)
	}

	return fmt.Sprintf(
		`<a href="/articles/%s/%s/">%s - %s</a><br>`,
		formatedTime,
		a.FileNameWithoutExtension(),
		formatedTime,
		a.Title,
	)
}

func (a *Article) FormatTime() string {
	return a.Timestamp.Format(timeformat)
}

func (a *Article) ToURL(fqdn string) string {
	return fmt.Sprintf("https://%s/articles/%s/%s", fqdn, a.FormatTime(), a.FileNameWithoutExtension())
}

func GenerateArticlePageHTML(a *Article, en bool) string {
	home := "/"
	if !en {
		home = "/ja"
	}

	return generateHTMLPage(fmt.Sprintf("%s | dtyler.io", a.Title), fmt.Sprintf(
		articlePageMD,
		home,
		a.Title,
		a.FormatTime(),
		strings.Join(a.ContentsMD, "\n"), home),
	)

}
