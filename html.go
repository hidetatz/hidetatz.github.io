package main

import (
	"fmt"
)

const (
	page = `
<!doctype html>
<html lang="en">
%s
%s
</html>
`

	head = `
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

	body = `
<body class="markdown-body">
%s

<script src="/syntax.js"></script>
<script>hljs.highlightAll();</script>
</body>
`
)

// body must be html
func generateHTMLPage(title, content string) string {
	return fmt.Sprintf(page, fmt.Sprintf(head, title), fmt.Sprintf(body, content))
}
