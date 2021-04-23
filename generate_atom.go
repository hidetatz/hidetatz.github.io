package main

import (
	"log"
	"time"

	"github.com/gorilla/feeds"
)

func genAtom(articles []*Article, count int, fqdn string) string {
	if count > len(articles) {
		count = len(articles)
	}

	name := "Hidetatsu Yaginuma"
	email := "deetyler@protonmail.com"
	feed := &feeds.Feed{
		Title:   "dtyler.io | Hidetatsy Yaginuma",
		Link:    &feeds.Link{Href: "https://dtyler.io"},
		Author:  &feeds.Author{Name: name, Email: email},
		Created: time.Now(),
		Items:   make([]*feeds.Item, count),
	}

	for i := 0; i < count; i++ {
		a := articles[i]
		feed.Items[i] = &feeds.Item{
			Title:       a.Title,
			Link:        &feeds.Link{Href: a.ToURL(fqdn)},
			Description: "The post first appeared on dtyler.io.",
			Author:      &feeds.Author{Name: name, Email: email},
			Created:     a.Timestamp,
		}

	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	return atom
}
