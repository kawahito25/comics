package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type comicInfo struct {
	title       string
	arrivalDate string
	shelfNumber string // 棚番号
}

type comicToRead struct {
	title  string
	toRead bool
}

type comics []*comicInfo

func getIsTarget(targetComicTitles []string, title string) bool {
	isTarget := false
	for _, target := range targetComicTitles {
		if strings.Contains(title, target) {
			isTarget = true
		}
	}
	return isTarget
}

func scrape(targetComicTitles []string) comics {
	collector := colly.NewCollector(
		colly.AllowedDomains(os.Getenv("ALLOWED_DOMAIN")),
	)
	detailCollector := collector.Clone()

	allComics := []comicToRead{}  // 今月入荷の全マンガ
	allArrivalDates := []string{} // 今月入荷の全マンガの入荷日
	targetShelfNums := []string{} // 読みたいマンガの棚番号

	collector.OnHTML(".arrival_detail", func(e *colly.HTMLElement) {
		title := e.ChildText(".title")
		detailLink := e.ChildAttr("a", "href")
		fmt.Println(title)

		c := comicToRead{
			title: title,
		}

		if getIsTarget(targetComicTitles, title) {
			c.toRead = true
			time.Sleep(2 * time.Second)
			detailCollector.Visit(e.Request.AbsoluteURL(detailLink))
		} else {
			c.toRead = false
		}

		allComics = append(allComics, c)
	})

	collector.OnHTML(".arrival_date", func(e *colly.HTMLElement) {
		allArrivalDates = append(allArrivalDates, e.Text)
	})

	collector.OnHTML("#list_header .pager", func(e *colly.HTMLElement) {
		e.ForEach("ul li", func(_ int, el *colly.HTMLElement) {
			if el.Text == "次>" {
				nextLink := el.ChildAttr("a", "href")
				time.Sleep(2 * time.Second)
				e.Request.Visit(nextLink) // 次のコミック一覧ページに移行
			}
		})
	})

	detailCollector.OnHTML(".last .more", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText("a[href]"))
		targetShelfNums = append(targetShelfNums, e.ChildText("a[href]"))
	})

	// Start scraping
	collector.Visit(os.Getenv("TARGET_URL"))

	comics := []*comicInfo{}
	for i, c := range allComics {
		if c.toRead {
			comics = append(comics, &comicInfo{
				title:       c.title,
				arrivalDate: allArrivalDates[i],
				shelfNumber: targetShelfNums[len(comics)],
			})
		}
	}

	return comics
}
