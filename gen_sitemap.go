package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/snabb/sitemap"
)

func genSiteMap(articles []*article, t time.Time, fqdn string) string {
	sm := sitemap.New()
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s", fqdn), LastMod: &t})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/ja/", fqdn), LastMod: &t})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/about/", fqdn), LastMod: &t})

	for _, a := range articles {
		sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/%s", fqdn, link(a)), LastMod: &a.timestamp})
	}

	buff := &bytes.Buffer{}
	sm.WriteTo(buff)
	return buff.String()
}
