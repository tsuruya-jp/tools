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

		e.ForEach(".design1 > tbody > tr:nth-child(1) > td:nth-child(4)", func(_ int, h *colly.HTMLElement) {
			caliber := h.Text
			caliber, _ = sjis_to_utf8(caliber)
			caliber = strings.TrimSpace(strings.TrimSuffix(caliber, "\n"))
			calibers = append(calibers, caliber)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(1) > td:nth-child(6)", func(_ int, h *colly.HTMLElement) {
			maker := h.Text
			maker, _ = sjis_to_utf8(maker)
			maker = strings.TrimSpace(strings.TrimSuffix(maker, "\n"))
			makers = append(makers, maker)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(2) > td:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			model := h.Text
			model, _ = sjis_to_utf8(model)
			model = strings.TrimSpace(strings.TrimSuffix(model, "\n"))
			models = append(models, model)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(2) > td:nth-child(4)", func(_ int, h *colly.HTMLElement) {
			gunHeight := h.Text
			gunHeight, _ = sjis_to_utf8(gunHeight)
			gunHeight = strings.TrimSpace(strings.TrimSuffix(gunHeight, "\n"))
			gunHeights = append(gunHeights, gunHeight)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(2) > td:nth-child(6)", func(_ int, h *colly.HTMLElement) {
			gunWeight := h.Text
			gunWeight, _ = sjis_to_utf8(gunWeight)
			gunWeight = strings.TrimSpace(strings.TrimSuffix(gunWeight, "\n"))
			gunWeights = append(gunWeights, gunWeight)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(3) > td:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			pull := h.Text
			pull, _ = sjis_to_utf8(pull)
			pull = strings.TrimSpace(strings.TrimSuffix(pull, "\n"))
			pulls = append(pulls, pull)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(3) > td:nth-child(4)", func(_ int, h *colly.HTMLElement) {
			condition := h.Text
			condition, _ = sjis_to_utf8(condition)
			condition = strings.TrimSpace(strings.TrimSuffix(condition, "\n"))
			conditions = append(conditions, condition)
		})

		e.ForEach(".design1 > tbody > tr:nth-child(3) > td:nth-child(6)", func(_ int, h *colly.HTMLElement) {
			remark := h.Text
			remark, _ = sjis_to_utf8(remark)
			remark = strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(remark, "\n"), "\n", " "))
			remarks = append(remarks, remark)
		})

		e.ForEach(".design1 > tbody > tr > td:nth-child(1)", func(_ int, h *colly.HTMLElement) {
			price := h.Text
			price, _ = sjis_to_utf8(price)
			price = strings.TrimSpace(strings.TrimSuffix(price, "\n"))
			prices = append(prices, price)
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
