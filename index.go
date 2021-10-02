package main

import (
	"fmt"
)

const indexPageMD = `
# dtyler.io

[About](/about) / [*Ja*](/ja) / [Input](/inputs) / [feed](/feed.xml)

---

%s
`

const jaIndexPageMD = `
# dtyler.io

[About](/about) / [*En*](/) / [Input](/inputs) / [feed](/feed_ja.xml)

---

%s
`

func generateIndexPageHTML(articlesList string) string {
	return generateHTMLPage("dtyler.io", fmt.Sprintf(indexPageMD, articlesList))
}

func generateJaIndexPageHTML(articlesList string) string {
	return generateHTMLPage("dtyler.io - ja", fmt.Sprintf(jaIndexPageMD, articlesList))
}
