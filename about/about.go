package about

import "github.com/gomarkdown/markdown"

const (
	tmpl = `
## About

Hidetatsu is a software engineer based on Tokyo, currently working for [Mercari, inc.](https://about.mercari.com/en/).

## Technical Skills

* Bachelor of Engineering degree (Computer science)
* Go - 4 Years+ (mainly used in my current job)
* Java - 5 Years+
* C++ - 6 Years+
* Designing Cloud platform architecture
  - Certification:
    - AWS Certified Solutions Architect - Professional (2017/12)
    - Google Cloud Platform Professional Cloud Architect (2019/12)
* Designing and building system architecture based on cloud platform

## Language Skills

* American English - Upper-intermediate
* Japanese - Native
`
)

type About struct {
	About string
}

func New() *About {
	return &About{About: tmpl}
}

func (a *About) ToHTML() string {
	return string(markdown.ToHTML([]byte(a.About), nil, nil))
}
