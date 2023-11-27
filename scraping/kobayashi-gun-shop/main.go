package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"	
)

const URL = "http://www.kobayashi-guns.co.jp/guns/real/guns_real.html"

func main() {
	c := colly.NewCollector()

	c.OnHTML("div[id=guns_real] > h2", func(e *colly.HTMLElement) {
		item := e.Text
		fmt.Println(item)
	})

	c.Visit(URL)
}
