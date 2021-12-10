package main

import "fmt"

const indexPageMD = `
# hidetatz.io

hidetatz.io is my personal website. The author Hidetatz (pronounced he-day-tatz) is a software engineer mainly focuses on system architecture, reliability, performance and observability based in Japan. I write code around infrastructure, database, transaction, concurrent programming and distributed systems. My code is available in [GitHub](https://github.com/hidetatz).

I [do fail](https://hidetatz.fail/).

---

## Projects

Give them a star!

### kubecolor (Go)

https://github.com/hidetatz/kubecolor

kubecolor is a CLI tool which wraps kubectl and colorizes the output for readability.
You can read my [blog article](https://hidetatz.medium.com/colorize-kubectl-output-by-kubecolor-2c222af3163a) about it.

### collection (Go)

https://github.com/hidetatz/collection

collection is a generics-aware Go library which provides collection data structures like [Java's one](https://docs.oracle.com/javase/8/docs/api/java/util/Collections.html).

### size-limited-queue (Go)

https://github.com/hidetatz/size-limited-queue

size-limited-queue is a blocking queue implementation. Internally sync.Cond is used and I made this repository to [describe](https://hidetatz.io/articles/2021/04/13/sync_cond/) how to use it.

---

## Articles

%s

Some articles are available in Japanese also.

%s

---

## Some other pages

* [/inputs](/inputs.html)
  - What I've read, listened, watched, etc.
* [/distsys](/distsys.html)
  - Distributed systems learning meterials (in Japanese)

---

If you want to send me any feedback about this website, you can submit it as GitHub issue [here](https://github.com/hidetatz/blog/issues/new).

---

[feed](/feed.xml)
`

func generateIndexPageHTML(articles []*article) string {
	enblogsList := ""
	jablogsList := ""
	for _, a := range articles {
		switch a.lang {
		case en:
			enblogsList += fmt.Sprintf("%s	- [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
		case ja:
			jablogsList += fmt.Sprintf("%s	- [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
		}
	}

	return generateHTMLPage("hidetatz.io", fmt.Sprintf(indexPageMD, enblogsList, jablogsList))
}

func link(a *article) string {
	formattedTime := a.timestamp.Format(timeformat)
	// if blog, the link should be external URL or internal link
	switch {
	case a.url == nil:
		return fmt.Sprintf("/articles/%s/%s", formattedTime, trimExtension(a.fileName))
	default:
		return a.url.String()
	}
}
