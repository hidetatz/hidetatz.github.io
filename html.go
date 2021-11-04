package main

import (
	"fmt"

	"github.com/russross/blackfriday/v2"
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
  <meta name="author" content="Hidetatz Yaginuma">
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

  <link href="/markdown.css" rel="stylesheet"></link>
  <link href="/syntax.css" rel="stylesheet"></link>
  <script type="text/javascript" async
    src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-MML-AM_CHTML">
  </script>
</head>
`

	Body = `
<body class="markdown-body">
%s

<hr>

<footer>
<p style="text-align:center">© 2017-2021 Hidetatz Yaginuma</p>
</footer>
<script src="/syntax.js"></script>
<script>hljs.highlightAll();</script>
</body>
`
)

func generateHTMLPage(title, contentsMarkdown string) string {
	head := fmt.Sprintf(Head, title)

	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	})
	renderer := &renderer{r}
	bodyHTML := string(blackfriday.Run([]byte(contentsMarkdown), blackfriday.WithRenderer(renderer)))
	body := fmt.Sprintf(Body, bodyHTML)

	return fmt.Sprintf(Page, head, body)
}
