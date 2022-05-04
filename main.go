package main

import (
	"fmt"
	"time"
	"strings"
	"os"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type comic struct {
	title       string
	arrivalDate string
	shelfNumber string // 棚番号
}

func main() {
	// 環境変数読み込み
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	c := colly.NewCollector(
		colly.AllowedDomains(os.Getenv("ALLOWED_DOMAIN")),
	)

	detailCollector := c.Clone()

	targetComics := []string{ // ToDo: Notion API を叩くようにする
		"僕のヒーローアカデミア",
		"キングダム",
	}

	comics := []comic{}
	arrivalDates := []string{}

	c.OnHTML(".arrival_date", func(e *colly.HTMLElement) {
		date := e.Text
		fmt.Println(date)
		arrivalDates = append(arrivalDates, date)
	})

	c.OnHTML(".arrival_detail", func(e *colly.HTMLElement) {
		title := e.ChildText(".title")
		detailLink := e.ChildAttr("a", "href")

		isTarget := false
		for _, target := range targetComics {
			if strings.Contains(title, target) {
				isTarget = true
			}
		}

		if isTarget {
			c := comic{
				title: title,
			}
			comics = append(comics, c)

			time.Sleep(1 * time.Second)
			detailCollector.Visit(e.Request.AbsoluteURL(detailLink))
		}

	})

	c.OnHTML("#list_header .pager", func(e *colly.HTMLElement) {
		// 次のページに移行
		e.ForEach("ul li", func(_ int, el *colly.HTMLElement) {
			if el.Text == "次>" {
				nextLink := el.ChildAttr("a", "href")
				time.Sleep(1 * time.Second)
				e.Request.Visit(nextLink)
			}
		})
	})

	detailCollector.OnHTML("div[class=more]", func(e *colly.HTMLElement) {
		fmt.Println("hello from detailCollector")
		fmt.Println(e.ChildText("a[href]"))
	})

	// Start scraping
	c.Visit(os.Getenv("TARGET_URL"))
}
