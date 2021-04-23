package main

import (
	"fmt"
)

const indexPageMD = `
# dtyler.io

[About](/about)

[*Ja*](/ja)

[feed](/feed.xml)

---

%s
`

const jaIndexPageMD = `
# dtyler.io

[About](/about)

[*En*](/)

[feed](/feed_ja.xml)

---

%s
`

func GenerateIndexPageHTML(articlesList string) string {
	return GenerateHTMLPage("dtyler.io", fmt.Sprintf(indexPageMD, articlesList))
}

func GenerateJaIndexPageHTML(articlesList string) string {
	return GenerateHTMLPage("dtyler.io - ja", fmt.Sprintf(jaIndexPageMD, articlesList))
}
