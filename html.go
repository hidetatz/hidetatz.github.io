package main

import (
	"fmt"

	"github.com/gomarkdown/markdown"
)

const (
	Page = `
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
  <meta name="viewport" content="width=device-width, initial-scale=1, minimal-ui">

  <style>
    body {
      box-sizing: border-box;
      min-width: 200px;
      max-width: 980px;
      margin: 0 auto;
      padding: 45px;
    }
  </style>

  <link href="https://cdn.jsdelivr.net/npm/github-markdown-css@3.0.1/github-markdown.min.css" rel="stylesheet"></link>
  <script type="text/javascript" async
    src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-MML-AM_CHTML">
  </script>
</head>
`

	Body = `
<body class="markdown-body">
%s

<footer>
<p style="text-aligh:center">© 2017-2021 Hidetatsu Yaginuma</p>
</footer>
</body>
`
)

func GenerateHTMLPage(title, contentsMarkdown string) string {
	head := fmt.Sprintf(Head, title)

	bodyHTML := string(markdown.ToHTML([]byte(contentsMarkdown), nil, nil))
	body := fmt.Sprintf(Body, bodyHTML)

	return fmt.Sprintf(Page, head, body)
}
