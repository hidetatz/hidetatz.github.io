package about

import (
	"github.com/yagi5/blog/html"
)

const (
	tmpl = `

## [dtyler.io](/)

# About

Hidetatsu is a software engineer, currently working for [Mercari, inc.](https://about.mercari.com/en/).

[GitHub](https://github.com/yagi5)
`
)

type About struct {
	About string
}

func New() *About {
	return &About{About: tmpl}
}

func (a *About) ToHTML() string {
	return html.Format("about", html.NewFromMarkdown(a.About))
}
