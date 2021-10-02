package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/feeds"
)

func genAtom(articles []*article, count int, fqdn string) string {
	if count > len(articles) {
		count = len(articles)
	}

	name := "Hidetatsu Yaginuma"
	email := "deetyler@protonmail.com"
	feed := &feeds.Feed{
		Title:   "dtyler.io | Hidetatsu Yaginuma",
		Link:    &feeds.Link{Href: "https://dtyler.io"},
		Author:  &feeds.Author{Name: name, Email: email},
		Created: time.Now(),
		Items:   make([]*feeds.Item, count),
	}

	for i := 0; i < count; i++ {
		a := articles[i]
		feed.Items[i] = &feeds.Item{
			Title:       a.title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s/%s", fqdn, link(a))},
			Description: "The post first appeared on dtyler.io.",
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
