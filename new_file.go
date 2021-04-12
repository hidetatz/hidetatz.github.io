package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

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
