package article

import (
	"fmt"
	"path/filepath"
	"time"
)

type Article struct {
	FileName   string
	Title      string
	Timestamp  time.Time
	ContentsMD *Contents
}

func (a *Article) Link() string {
	return fmt.Sprintf(
		`<a href="/articles/%s/%s/">%s - %s</a><br>`,
		a.YMD(),
		a.FileName,
		a.Title,
		a.Timestamp.Format("2006/01/02"),
	)
}

func (a *Article) YMD() string {
	month := a.Timestamp.Month()
	smonth := fmt.Sprintf("%d", month)
	if month < 10 {
		smonth = fmt.Sprintf("0%d", month)
	}
	return fmt.Sprintf(
		"%d/%s/%d",
		a.Timestamp.Year(),
		smonth,
		a.Timestamp.Day(),
	)
}

func (a *Article) FileNameWithoutExtension() string {
	return filepath.Base(a.FileName[:len(a.FileName)-len(filepath.Ext(a.FileName))])
}
