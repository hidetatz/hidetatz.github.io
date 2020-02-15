package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yagi5/blog/about"
	"github.com/yagi5/blog/article"
	"github.com/yagi5/blog/cname"
	"github.com/yagi5/blog/css"
	"github.com/yagi5/blog/index"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd := os.Args[1]
	switch {
	case cmd == "new":
		if len(os.Args) != 3 {
			log.Fatalf("title is empty")
		}

		newArticle(os.Args[2])

	case cmd == "gen":
		removeAllFiles("./docs/")
		gen()
	}
}

func newArticle(title string) {
	content := fmt.Sprintf(`%s---%s`, title, time.Now().Format("2006-01-02 15:04:05"))
	write(content, fmt.Sprintf("./data/articles/%s.md", title))
	fmt.Printf("file created: %s\n", fmt.Sprintf("./data/articles/%s.md", title))
}

func gen() {
	cname := cname.New()
	write(cname, "./docs/CNAME")

	css := css.New()
	write(css, "./docs/markdown.css")

	articles := article.NewArticles("./data/articles")
	articleList := articles.List()

	for _, a := range articles {
		write(a.ToHTML(), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.YMD(), a.FileNameWithoutExtension()))
	}

	articlesJA := article.NewArticles("./data/articles/ja")
	articlesJAList := articlesJA.ListJA()

	for _, a := range articlesJA {
		write(a.ToHTML(), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.YMD(), a.FileNameWithoutExtension()))
	}

	idx := index.New(articleList)
	idxJA := index.NewJA(articlesJAList)

	write(idx, "./docs/index.html")
	write(idxJA, "./docs/ja/index.html")

	about := about.New()
	write(about.ToHTML(), "./docs/about/index.html")
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
