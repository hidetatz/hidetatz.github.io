package main

import (
	"fmt"
)

const indexPageMD = `
# dtyler.io

[About](/about)

[Ja](/ja)

---

%s
`

const jaIndexPageMD = `
# dtyler.io

[About](/about)

[En](/)

---

%s
`

func GenerateIndexPageHTML(articlesList string) string {
	return GenerateHTMLPage("dtyler.io", fmt.Sprintf(indexPageMD, articlesList))
}

func GenerateJaIndexPageHTML(articlesList string) string {
	return GenerateHTMLPage("dtyler.io - ja", fmt.Sprintf(jaIndexPageMD, articlesList))
}
