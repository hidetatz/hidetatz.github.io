package main

import "fmt"

const indexPageMD = `
# dtyler.io

dtyler.io is my personal blog. The author Hidetatsu is a software engineer mainly focuses on reliability, performance, observability and developer experience based in Japan. Hidetatsu does infrastructure, database, transaction, concurrent programming and distributed systems. My code is available in [GitHub](https://github.com/dty1er).

---

## Articles

%s

Some articles are available in Japanese also.

%s

---

## Inputs

What I've read/learned, with some thoughts.

%s

---

[feed](/feed.xml)
`

func generateIndexPageHTML(articles []*article) string {
	enblogsList := ""
	jablogsList := ""
	inputsList := ""
	for _, a := range articles {
		switch a.typ {
		case inputType:
			inputsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
		case blogType:
			switch a.lang {
			case en:
				enblogsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
			case ja:
				jablogsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
			}
		}
	}

	return generateHTMLPage("dtyler.io", fmt.Sprintf(indexPageMD, enblogsList, jablogsList, inputsList))
}

func link(a *article) string {
	formattedTime := a.timestamp.Format(timeformat)
	switch a.typ {
	case blogType:
		// if blog, the link should be external URL or internal link
		switch {
		case a.url == nil:
			return fmt.Sprintf("/articles/%s/%s", formattedTime, trimExtension(a.fileName))
		default:
			return a.url.String()
		}
	default:
		// else, it is input. return internal link
		return fmt.Sprintf("/articles/%s/%s", formattedTime, trimExtension(a.fileName))
	}
}
