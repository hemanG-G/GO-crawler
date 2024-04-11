package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
}
type movie struct {
	Title string
	Year  string
}

func main() {
	month := flag.Int("month", 1, "Month To Fetch Birthday for")
	day := flag.Int("day", 1, "Day to fetch Birthday for")
	flag.Parse()

	crawl(*month, *day) // pointers cuz we parse these from user
}

func crawl(month int, day int) {

	c := colly.NewCollector( // 1st colly instance to go to profiles
		colly.AllowedDomains("imdb.com", "www.imdb.com", "www.web.archive.org", "web.archive.org"),
	)
	infoCollector := c.Clone() // 2nd colly instance that goes into each of the Profiles

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	// HTML Element from IMBD page inspect
	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) { // going to next page functionality to the infoCollector
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)

	})

	infoCollector.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		tempProfile := star{}
		tempProfile.Name = e.ChildText("hl.header > span.itemprop")
		tempProfile.Photo = e.ChildAttr("#name-poster", "src")
		tempProfile.JobTitle = e.ChildText("#name-job-categories > a > span.itemprop")
		tempProfile.BirthDate = e.ChildAttr("#name-born-info time", "datetime")
		tempProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

		e.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = kf.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = kf.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tempProfile.TopMovies = append(tempProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tempProfile, "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Profile URL ", r.URL.String())
	})

	// wayback machine Link  , as HTML is updated now
	startUrl := fmt.Sprintf("https://web.archive.org/web/20210125231114/https://www.imdb.com/search/name/?birth_monthday=01-01")
	c.Visit(startUrl)

}
