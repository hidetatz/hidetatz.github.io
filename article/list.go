package article

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const tmpl = `
<a href="https://medium.com/@yagi5/lets-study-distributed-systems-4-leader-election-78a083981321">Let’s study distributed systems — 4. Leader election</a><br>
<a href="https://medium.com/@yagi5/lets-study-distributed-systems-3-distributed-snapshots-4bd8eba07988">Let’s study distributed systems — 3. Distributed snapshot</a><br>
<a href="https://medium.com/@yagi5/lets-study-distributed-systems-2-clock-67ffcf98a645">Let’s study distributed systems — 2. Clock<a><br>
<a href="https://medium.com/@yagi5/lets-study-distributed-systems-1-introduction-e149f2157253">Let’s study distributed systems — 1. Introduction</a><br>
%s
`

type Articles []*Article

func (as Articles) List() string {
	sb := strings.Builder{}
	for _, article := range as {
		sb.WriteString(article.Link() + "\n")
	}

	return fmt.Sprintf(tmpl, sb.String())
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
				article.Title, article.Timestamp = parseFirstLine(file.Name(), scanner.Text())
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

func parseFirstLine(fileName, firstLine string) (string, time.Time) {
	// must be formatted
	// title of article---2019-04-01 12:30:00
	splitted := strings.Split(firstLine, "---")
	if len(splitted) != 2 {
		log.Fatalf("Invalid format of article: %s, %s", fileName, firstLine)
	}

	t, err := time.Parse("2006-01-02 15:04:05", splitted[1])
	if err != nil {
		log.Fatalf("Invalid time format: '%s' in %s", splitted[1], fileName)
	}
	return splitted[0], t
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
