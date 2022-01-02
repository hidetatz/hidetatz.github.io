package main

import (
	"fmt"
)

const indexPageMD = `
# hidetatz.io

hidetatz.io is my personal website. The author Hidetatz (pronounced he-day-tatz) is a software engineer mainly focuses on system architecture, reliability, performance and observability based in Japan. I write code around infrastructure, database, transaction, concurrent programming and distributed systems. My code is available in [GitHub](https://github.com/hidetatz).

If you want to send me any feedback or questions about this website/article, you can submit it as GitHub issue [here](https://github.com/hidetatz/blog/issues/new).

Atom/RSS feed is found [here](/feed.xml).

I [do fail](https://hidetatz.fail/).

---

## Projects

If you love it, give a star!

* [kubecolor](https://github.com/hidetatz/kubecolor) (Go)
  - kubecolor is a CLI tool which wraps kubectl and colorizes the output for readability.
  - You can read my [blog article](https://hidetatz.medium.com/colorize-kubectl-output-by-kubecolor-2c222af3163a) about it.
* [collection](https://github.com/hidetatz/collection) (Go)
  - collection is a generics-aware Go library which provides collection data structures like [Java's one](https://docs.oracle.com/javase/8/docs/api/java/util/Collections.html).
* [size-limited-queue](https://github.com/hidetatz/size-limited-queue) (Go)
  - size-limited-queue is a blocking queue implementation. Internally sync.Cond is used and I made this repository to [describe](https://hidetatz.io/articles/2021/04/13/sync_cond/) how to use it.

---

## Articles

%s

Some articles are available in Japanese also.

%s

---

## Other writings

* [/distsys](/distsys.html)
  - Distributed systems learning meterials (in Japanese)

---

© 2022 Hidetatz Yaginuma. Unless otherwise noted, these posts are made available under a [Creative Commons Attribution License](https://creativecommons.org/licenses/by/4.0/).
`

func generateIndexPageHTML(articles []*article) string {
	enblogsList := ""
	jablogsList := ""
	for _, a := range articles {
		switch a.lang {
		case en:
			enblogsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, linkToArticle(a))
		case ja:
			jablogsList += fmt.Sprintf("%s - [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, linkToArticle(a))
		}
	}

	contentsHTML := toHTML(fmt.Sprintf(indexPageMD, enblogsList, jablogsList))
	return generateHTMLPage("hidetatz.io", contentsHTML)
}
