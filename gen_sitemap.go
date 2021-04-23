package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/snabb/sitemap"
)

func genSiteMap(articles []*Article, fqdn string) string {
	sm := sitemap.New()
	now := time.Now()
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s", fqdn), LastMod: &now})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/ja/", fqdn), LastMod: &now})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/about/", fqdn), LastMod: &now})

	for _, a := range articles {
		loc := fmt.Sprintf("https://%s/articles/%s/%s/", fqdn, a.FormatTime(), a.FileNameWithoutExtension())
		sm.Add(&sitemap.URL{Loc: loc, LastMod: &a.Timestamp})
	}

	buff := &bytes.Buffer{}
	sm.WriteTo(buff)
	return buff.String()
}
