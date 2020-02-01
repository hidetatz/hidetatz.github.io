package about

import "github.com/gomarkdown/markdown"

const (
	tmpl = `
## About

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
	ret := string(markdown.ToHTML([]byte(a.About), nil, nil))
	return `<link href="/markdown.css" rel="stylesheet"></link>` + "\n" + ret
}
