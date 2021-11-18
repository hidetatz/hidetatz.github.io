package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/feeds"
)

func genAtom(articles []*article, t time.Time, count int, fqdn string) string {
	if count > len(articles) {
		count = len(articles)
	}

	name := "Hidetatz Yaginuma"
	email := "hidetatz@gmail.com"
	feed := &feeds.Feed{
		Title:   fmt.Sprintf("hidetatz.io | %s", name),
		Link:    &feeds.Link{Href: "https://dtyler.io"},
		Author:  &feeds.Author{Name: name, Email: email},
		Created: t,
		Items:   make([]*feeds.Item, count),
	}

	for i := 0; i < count; i++ {
		a := articles[i]
		feed.Items[i] = &feeds.Item{
			Title:       a.title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s/%s", fqdn, link(a))},
			Description: "The post first appeared on hidetatz.io.",
			Author:      &feeds.Author{Name: name, Email: email},
			Created:     a.timestamp,
		}

	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	return atom
}
