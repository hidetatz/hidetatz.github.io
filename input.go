package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"

	xhtml "golang.org/x/net/html"
)

type inputArticle struct {
	title      string
	url        *url.URL
	timestamp  time.Time
	contentsMD []string
	FileName   string
}

func (i *inputArticle) FormatTime() string {
	return i.timestamp.Format(timeformat)
}

func (i *inputArticle) FileNameWithoutExtension() string {
	return filepath.Base(i.FileName[:len(i.FileName)-len(filepath.Ext(i.FileName))])
}

func readInputs(dir string) []*inputArticle {
	files := readFiles(dir)

	inputs := []*inputArticle{}
	for _, file := range files {
		scanner := bufio.NewScanner(file)
		firstLine := true
		lines := []string{}
		input := &inputArticle{}
		for scanner.Scan() {
			if firstLine {
				if err := decodeInputFirstLine(input, scanner.Text()); err != nil {
					log.Fatalf("failed to decode the article first line: %s, %s", file.Name(), err)
				}
				firstLine = false
				continue
			}
			lines = append(lines, scanner.Text())
		}
		input.contentsMD = lines
		input.FileName = file.Name()
		inputs = append(inputs, input)
	}

	sort.Slice(inputs, func(i, j int) bool { return inputs[i].timestamp.After(inputs[j].timestamp) })
	return inputs
}

func decodeInputFirstLine(input *inputArticle, firstLine string) error {
	splitted := strings.Split(firstLine, "---")
	if len(splitted) != 2 {
		return fmt.Errorf("invalid format of article first line: %s", firstLine)
	}

	u, err := url.Parse(splitted[0])
	if err != nil {
		return fmt.Errorf("invalid url: '%s'", splitted[0])
	}

	input.url = u

	t, err := time.Parse("2006-01-02 15:04:05", splitted[1])
	if err != nil {
		return fmt.Errorf("invalid time format: '%s'", splitted[1])
	}

	input.timestamp = t

	title, ok, err := getHTMLTitle(u)
	if err != nil {
		return fmt.Errorf("failed to get HTML title", splitted[0])
	}

	if !ok {
		input.title = u.String()
		return nil
	}

	input.title = title

	return nil
}

func isTitleElement(n *xhtml.Node) bool {
	return n.Type == xhtml.ElementNode && n.Data == "title"
}

func traverse(n *xhtml.Node) (string, bool) {
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

func getHTMLTitle(u *url.URL) (string, bool, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return "", false, fmt.Errorf("failed to call %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	doc, err := xhtml.Parse(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse %s: %w", u.String(), err)
	}

	title, ok := traverse(doc)
	return title, ok, nil
}

func listInputsHref(inputs []*inputArticle) []string {
	ret := make([]string, len(inputs))
	for i, input := range inputs {
		ret[i] = fmt.Sprintf(
			`<a href="/inputs/%s/%s/">%s - %s</a><br>`,
			input.FormatTime(),
			input.FileNameWithoutExtension(),
			input.FormatTime(),
			input.title,
		)
	}

	return ret
}

func generateInputPageHTML(i *inputArticle) string {
	home := "/"

	return generateHTMLPage(fmt.Sprintf("%s | dtyler.io", i.title), fmt.Sprintf(
		articlePageMD,
		home,
		i.title,
		i.FormatTime(),
		strings.Join(i.contentsMD, "\n"), home),
	)

}
