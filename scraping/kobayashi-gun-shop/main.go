package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const URL = "http://www.kobayashi-guns.co.jp/guns/real/guns_real.html"

func sjis_to_utf8(str string) (string, error) {
        ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
        if err != nil {
                return "", err
        }
        return string(ret), err
}

func main() {
	c := colly.NewCollector()

	records := make([][]string, 0)
	titles := make([]string, 0)
	kinds := make([]string, 0)
	i := 0
	c.OnHTML("div[id=guns_real]", func(e *colly.HTMLElement) {
		e.ForEach("h2", func(_ int, h *colly.HTMLElement) {
			i++
			title := h.Text
			title, _ = sjis_to_utf8(title)
			titles = append(titles, title)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(1) > td:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			kind := h.Text
			kind, _ = sjis_to_utf8(kind)
			kind = strings.TrimSpace(strings.TrimSuffix(kind, "\n"))
			kinds = append(kinds, kind)
		})

		records = append(records, titles)
		records = append(records, kinds)

	})

	c.Visit(URL)
	c.Wait()

	fmt.Println(i)

	file, err := os.Create("sample.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.WriteAll(records)
	w.Flush()
}
