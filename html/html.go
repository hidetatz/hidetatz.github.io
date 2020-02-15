package html

import (
	"fmt"

	"github.com/gomarkdown/markdown"
)

const (
	HTML = `
<!doctype html>
<html lang="en">
%s
%s
</html>
`

	Head = `
<head>
  <meta charset="utf-8">

  <title>%s</title>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="author" content="Hidetatsu Yaginuma">

  <link href="/markdown.css" rel="stylesheet"></link>
  <script type="text/javascript" async
    src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-MML-AM_CHTML">
  </script>
</head>
`

	Body = `
<body>
%s
</body>
`
)

func NewBody(contents string) string {
	return fmt.Sprintf(Body, contents)
}

func NewHead(title string) string {
	return fmt.Sprintf(Head, title)
}

func Format(title, contents string) string {
	head := NewHead(title)
	body := NewBody(contents)
	return fmt.Sprintf(HTML, head, body)
}

func NewFromMarkdown(md string) string {
	return string(markdown.ToHTML([]byte(md), nil, nil))
}
