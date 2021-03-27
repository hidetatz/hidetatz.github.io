package main

import (
	"fmt"
)

// indexPage contains title, link to the about page,
// link to the another locale(en->ja, ja->en) index page, and a list of articles url.

// const indexPage = `
// <a href="/"><h1>dtyler.io</h1></a>
// <a href="/about">About</a><br>
// <a href="/ja">Japanese articles</a><br>

// <hr>

// %s
// `

// const jaIndexPage = `
// <a href="/"><h1>dtyler.io</h1></a>
// <a href="/about">About</a><br>
// <a href="/">English articles</a><br>

// <hr>

// %s
// `

const indexPageMD = `
## [dtyler.io](/)

---

[About](/about)

[Ja](/ja)

---

%s
`

const jaIndexPageMD = `
## [dtyler.io](/)

---

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
