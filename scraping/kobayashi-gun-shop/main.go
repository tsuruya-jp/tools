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
	c.OnHTML("div[id=guns_real]", func(e *colly.HTMLElement) {
		e.ForEach("h2", func(_ int, h *colly.HTMLElement) {
			title := h.Text
			title, _ = sjis_to_utf8(title)
			titles = append(titles, title)
		})

		records = append(records, titles)

	})

	c.Visit(URL)
	c.Wait()

	fmt.Printf("%#v", records)

	file, err := os.Create("sample.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.WriteAll(records)
	w.Flush()
}
