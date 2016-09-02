package main

import (
	"fmt"
	"net/http"

	"github.com/goshinobi/client"
)

func main() {
	c := client.New(true)
	fmt.Println(c)
	c.Add(&http.Client{}, nil)
	fmt.Println(c)
	c.Add(&http.Client{}, nil)
	fmt.Println(c)
	fmt.Println(c.GetIPs())
	fmt.Println(c.Get("http://google.co.jp"))
}
