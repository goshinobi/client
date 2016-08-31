package main

import (
	"fmt"
	"net/http"

	"github.com/goshinobi/client"
)

func main() {
	c := client.New()
	fmt.Println(c)
	c.Add(&http.Client{}, 5)
	fmt.Println(c)
	c.Add(&http.Client{}, 5)
	fmt.Println(c)
}
