package main

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding"
	"log"
	"strings"
	"github.com/bartekbp/sportpl/crawler"
)

func Parse(ctx *gocrawl.URLContext, doc *goquery.Document, decoder *encoding.Decoder) crawler.Page {
	page := crawler.Page{Url: ctx.NormalizedURL().String()}
	path := ctx.NormalizedURL().Path
	pathParts := strings.SplitAfter(path, "/")
	if len(pathParts) != 3 {
		return page
	}

	ids := strings.Split(pathParts[2], ",")
	if len(ids) != 4 {
		return page
	}

	topic := ids[1]
	var comments []crawler.Comment
	doc.Find(".comment-body").Parent().Each(func(_ int, selection *goquery.Selection) {
		id := selection.AttrOr("data-id", "")
		text := selection.Find(".inner").ChildrenFiltered("p").Text()

		decoded, err := decoder.String(text)
		if err != nil {
			log.Fatal(err)
		}

		comments = append(comments, crawler.Comment{Id: id, Topic: topic, Text: strings.Trim(decoded, "\n\r\t ")})
	})

	page.Comments = comments
	return page
}
