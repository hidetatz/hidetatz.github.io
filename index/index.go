package index

import (
	"strings"

	"github.com/yagi5/blog/html"
)

const body = `
<a href="/"><h1>dtyler.io</h1></a>
<a href="/about">About</a><br>
${another_locale}

<hr>

${articles_list}
`

func New(articleList string) string {
	contents := strings.Replace(body, "${another_locale}", `<a href="/ja">Ja</a><br>`, -1)
	contents = replaceArticleList(contents, articleList)

	return html.Format("dtyler.io", contents)
}

func NewJA(articleList string) string {
	contents := strings.Replace(body, "${another_locale}", `<a href="/">En</a><br>`, -1)
	contents = replaceArticleList(contents, articleList)

	return html.Format("dtyler.io - ja", contents)
}

func replaceArticleList(orig, articleList string) string {
	return strings.Replace(orig, "${articles_list}", articleList, -1)
}
