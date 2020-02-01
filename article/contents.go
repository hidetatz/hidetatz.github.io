package article

import (
	"strings"

	"github.com/gomarkdown/markdown"
)

type Contents struct {
	lines []string
}

func NewContents(lines []string) *Contents {
	return &Contents{lines: lines}
}

func (c *Contents) ToHTML() string {
	cts := strings.Join(c.lines, "\n")
	ret := string(markdown.ToHTML([]byte(cts), nil, nil))
	return `<link href="/markdown.css" rel="stylesheet"></link>` + "\n" + ret
}
