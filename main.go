package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/fsnotify/fsnotify"
)

const (
	cname = "dtyler.io"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	g := flag.Bool("gen", false, "generate static site from markdown articles")
	n := flag.Bool("new", false, "generate a new article file template")

	flag.Parse()

	switch {
	case *g:
		removeAllFiles("./docs/")
		gen()
	case *n:
		newFile()
	default:
		runServer()
	}
}

func newFile() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("en or ja?")
	lang, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	lang = strings.TrimSpace(lang)
	if lang != "en" && lang != "ja" {
		log.Fatal("lang must either en or ja")
	}

	fmt.Println("filename?")
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		log.Fatal("file must not be empty")
	}

	fmt.Println("url? (just hit enter if not an external URL)")
	url, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	url = strings.TrimSpace(url)
	templateContent := fmt.Sprintf("title of article here---%s", time.Now().Format("2006-01-02 15:04:05"))
	if url != "" {
		templateContent += fmt.Sprintf("---%s", url)
	}

	if lang == "en" {
		write(templateContent, fmt.Sprintf("./data/articles/%s.md", filename))
	} else {
		write(templateContent, fmt.Sprintf("./data/articles/ja/%s.md", filename))
	}
}

func runServer() {
	ctx := context.Background()

	port := "8080"
	directory := "./docs"
	ip := getLocalIP()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	server := &http.Server{Addr: ":8080", Handler: http.FileServer(http.Dir(directory))}

	go func() {
		server.ListenAndServe()
	}()

	fmt.Printf("Serving at %s:%s (the address is automatically copied into the clipboard)\n", ip, port)
	clipboard.WriteAll(fmt.Sprintf("%s:%s", ip, port))

	// live reload with fsnotify
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}

				fmt.Println("files edited. hot reloading...")
				server.Shutdown(ctx)
				time.Sleep(time.Second) // wait for shutdown

				removeAllFiles("./docs/")
				gen()
				server = &http.Server{Addr: ":8080", Handler: http.FileServer(http.Dir(directory))}
				// don't take care of goroutine cancel since this is not a long-running script and
				// want to make code look easy
				go func() {
					server.ListenAndServe()
				}()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	watcher.Add("./data/articles/ja")
	watcher.Add("./data/articles/")
	<-done
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func gen() {
	write(cname, "./docs/CNAME") // required for GitHub pages
	write(css, "./docs/markdown.css")
	cp("favicon.ico", "./docs/favicon.ico")
	write(GenerateHTMLPage("about", About), "./docs/about/index.html")

	articles := ReadArticles("./data/articles")
	articleList := ListArticlesHref(articles)

	for _, a := range articles {
		// if url != nil, no need to generate the page because it is linked from nowhere
		if a.URL == nil {
			write(GenerateArticlePageHTML(a), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.FormatTime(), a.FileNameWithoutExtension()))
		}
	}

	articlesJA := ReadArticles("./data/articles/ja")
	articlesJAList := ListArticlesHref(articlesJA)

	for _, a := range articlesJA {
		if a.URL == nil {
			write(GenerateArticlePageHTML(a), fmt.Sprintf("./docs/articles/%s/%s/index.html", a.FormatTime(), a.FileNameWithoutExtension()))
		}
	}

	idx := GenerateIndexPageHTML(strings.Join(articleList, "\n"))
	idxJA := GenerateJaIndexPageHTML(strings.Join(articlesJAList, "\n"))

	write(idx, "./docs/index.html")
	write(idxJA, "./docs/ja/index.html")

}

func cp(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal(err)
	}
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
