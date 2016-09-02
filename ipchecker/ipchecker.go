package ipchecker

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

var (
	buffer []func(*http.Client) string
	p      int
)

func init() {
	RegistChecker(cman)
}

func RegistChecker(f func(*http.Client) string) {
	buffer = append(buffer, f)
}

func Check(c *http.Client) string {
	result := buffer[p](c)
	p++
	if p >= len(buffer) {
		p = 0
	}
	return result
}

var cman = func(c *http.Client) string {
	url := "https://www.cman.jp/network/support/go_access.cgi"

	resp, err := c.Get(url)
	if err != nil {
		return ""
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return ""
	}

	return doc.Find(".outIp").Text()
}
