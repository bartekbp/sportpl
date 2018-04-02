package main

import (
	"fmt"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/bartekbp/sportpl/crawler"
	"log"
	"regexp"
	"strings"
	"time"
)

func processPage(page crawler.Page) {
	if len(page.Comments) == 0 {
		return
	}

	for _, comment := range page.Comments {
		fakeComment := strings.Contains(comment.Text, "sprawdzisz swoje IQ")

		if fakeComment {
			log.Println("Fake comment on:", page.Url)
			baseUrl := "http://www.sport.pl/fix/cms/opinions/opinions-action.jsp?"
			fullUrl := fmt.Sprintf(baseUrl+"action=trashVote&id=%s&dzialXx=%s&jspXx=", comment.Id, comment.Topic)
			log.Println(fullUrl)
		}
	}
}

func main() {
	sportPlDomain := regexp.MustCompile("http://(www.)?sport.pl(/.*)?")
	extender := &crawler.CommentExtender{Pages: make(chan crawler.Page, 100), Domain: sportPlDomain, Parse: Parse}
	opts := gocrawl.NewOptions(extender)

	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0)"
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = 5
	opts.SameHostOnly = false

	c := gocrawl.NewCrawlerWithOptions(opts)
	doneCrawling := make(chan bool, 1)

	go func() {
		c.Run("http://www.sport.pl")
		doneCrawling <- true
	}()

	crawler.WaitForPages(extender.Pages, doneCrawling, processPage)
}
