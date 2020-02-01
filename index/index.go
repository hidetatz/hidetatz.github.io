package index

import "strings"

const tmpl = `<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

  <title>dtyler.io</title>
  <meta name="description" content="dtyler.io">
  <meta name="author" content="SitePoint">

  <link href="/markdown.css" rel="stylesheet"></link>

  <!--[if lt IE 9]>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/html5shiv/3.7.3/html5shiv.js"></script>
  <![endif]-->
</head>

<body>
    <h1>dtyler.io</h1>

    <a href="/about">About</a><br>
	${another_locale}

    ${articles_list}
</body>
</html>
`

func New(articleList string) string {
	ret := strings.Replace(tmpl, "${another_locale}", `<a href="/ja">Ja</a><br>`, -1)
	return replaceArticleList(ret, articleList)
}

func NewJA(articleList string) string {
	ret := strings.Replace(tmpl, "${another_locale}", `<a href="/">En</a><br>`, -1)
	return replaceArticleList(ret, articleList)
}

func replaceArticleList(orig, articleList string) string {
	return strings.Replace(orig, "${articles_list}", articleList, -1)
}
