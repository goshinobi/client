package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/goshinobi/client"
	"github.com/goshinobi/ipchecker"
)

func init() {
	ipchecker.RegistChecker(myIPChecker)
}

func myIPChecker(c *http.Client) string {
	result := ""
	url := "http://www.ugtop.com/spill.shtml"
	resp, err := c.Get(url)
	if err != nil {
		return ""
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	doc.Find("table").Find("tbody").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				if i == 2 {
					result = s.Find("font").Text()
				}
			})
		}
	})

	return result
}

func main() {
	c := client.NewClient(client.ClientCfg{}, client.ClientCfg{}, client.ClientCfg{})
	fmt.Println(c.Len())
	torClient := client.NewClientTor(10)
	c.Add(torClient)

	fmt.Println(c.Len())
}
