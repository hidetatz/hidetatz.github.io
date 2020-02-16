package article

import (
	"bufio"
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

type Articles []*Article

func (as Articles) List() string {
	sb := strings.Builder{}
	for _, article := range as {
		sb.WriteString(article.Link() + "\n")
	}

	return sb.String()
}

func (as Articles) ListJA() string {
	sb := strings.Builder{}
	for _, article := range as {
		sb.WriteString(article.Link() + "\n")
	}

	return sb.String()
}

func NewArticles(dir string) Articles {
	files := articleFiles(dir)

	articles := []*Article{}
	for _, file := range files {
		scanner := bufio.NewScanner(file)
		firstLine := true
		lines := []string{}
		article := &Article{}
		for scanner.Scan() {
			if firstLine {
				article.Title, article.Timestamp, article.URL = parseFirstLine(file.Name(), scanner.Text())
				firstLine = false
				continue
			}
			lines = append(lines, scanner.Text())
		}
		article.FileName = filepath.Base(file.Name())
		article.ContentsMD = NewContents(lines)
		articles = append(articles, article)
	}

	sort.Slice(articles, func(i, j int) bool { return articles[i].Timestamp.After(articles[j].Timestamp) })
	return articles
}

func parseFirstLine(fileName, firstLine string) (string, time.Time, *url.URL) {
	// must be formatted
	// title of article---2019-04-01 12:30:00
	// or
	// title of article---2019-04-01 12:30:00---https://url.to.article
	splitted := strings.Split(firstLine, "---")
	if len(splitted) != 2 && len(splitted) != 3 {
		log.Fatalf("Invalid format of article: %s, %s", fileName, firstLine)
	}

	t, err := time.Parse("2006-01-02 15:04:05", splitted[1])
	if err != nil {
		log.Fatalf("Invalid time format: '%s' in %s", splitted[1], fileName)
	}

	if len(splitted) == 3 {
		u, err := url.Parse(splitted[2])
		if err != nil {
			log.Fatalf("failed to parse URLt: '%s' in %s", splitted[2], fileName)
		}
		return splitted[0], t, u
	}

	return splitted[0], t, nil
}

func articleFiles(dir string) []*os.File {
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
