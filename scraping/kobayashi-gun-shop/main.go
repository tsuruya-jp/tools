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

func sjis_to_utf8(str string) string {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		panic(err)	
	}
	return string(ret)
}

func transpose(slice [][]string) [][]string {
    xl := len(slice[0])
    yl := len(slice)
    result := make([][]string, xl)
    for i := range result {
        result[i] = make([]string, yl)
    }
    for i := 0; i < xl; i++ {
        for j := 0; j < yl; j++ {
            result[i][j] = slice[j][i]
        }
    }
    return result
}

func main() {
	c := colly.NewCollector()

	records := make([][]string, 0)
	titles := make([]string, 0)
	kinds := make([]string, 0)
	calibers := make([]string, 0)
	makers := make([]string, 0)
	models := make([]string, 0)
	gunHeights := make([]string, 0)
	gunWeights := make([]string, 0)
	pulls := make([]string, 0)
	conditions := make([]string, 0)
	remarks := make([]string, 0)
	prices := make([]string, 0)
	i := 0
	c.OnHTML("div[id=guns_real]", func(e *colly.HTMLElement) {
		e.ForEach("h2", func(_ int, h *colly.HTMLElement) {
			i++
			titles = append(titles, sjis_to_utf8(h.Text))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(1) > td:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			kinds = append(kinds, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(1) > td:nth-child(4)", func(_ int, h *colly.HTMLElement) {
			calibers = append(calibers, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(1) > td:nth-child(6)", func(_ int, h *colly.HTMLElement) {
			makers = append(makers, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(2) > td:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			models = append(models, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(2) > td:nth-child(4)", func(_ int, h *colly.HTMLElement) {
			gunHeights = append(gunHeights, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(2) > td:nth-child(6)", func(_ int, h *colly.HTMLElement) {
			gunWeights = append(gunWeights, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(3) > td:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			pulls = append(pulls, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(3) > td:nth-child(4)", func(_ int, h *colly.HTMLElement) {
			conditions = append(conditions, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		e.ForEach(".design1 > tbody > tr:nth-child(3) > td:nth-child(6)", func(_ int, h *colly.HTMLElement) {
			remarks = append(remarks, strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n"), "\n", " ")))
		})

		e.ForEach(".design1 > tbody > tr > td:nth-child(1)", func(_ int, h *colly.HTMLElement) {
			prices = append(prices, strings.TrimSpace(strings.TrimSuffix(sjis_to_utf8(h.Text), "\n")))
		})

		records = append(records, titles)
		records = append(records, kinds)
		records = append(records, calibers)
		records = append(records, makers)
		records = append(records, models)
		records = append(records, gunHeights)
		records = append(records, gunWeights)
		records = append(records, pulls)
		records = append(records, conditions)
		records = append(records, remarks)
		records = append(records, prices)

	})

	c.Visit(URL)
	c.Wait()

	records = transpose(records)

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
