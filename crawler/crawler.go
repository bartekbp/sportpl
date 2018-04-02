package crawler

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"log"
	"net/http"
	"regexp"
)

type Comment struct {
	Id    string
	Topic string
	Text  string
}

type Page struct {
	Url      string
	Comments []Comment
}

type CommentExtender struct {
	gocrawl.DefaultExtender
	Pages  chan Page
	Domain *regexp.Regexp
	Parse  func(ctx *gocrawl.URLContext, doc *goquery.Document, decoder *encoding.Decoder) Page
}

func (x *CommentExtender) Visit(ctx *gocrawl.URLContext, response *http.Response, doc *goquery.Document) (interface{}, bool) {
	statusCode := response.StatusCode
	if statusCode < 200 || statusCode > 300 {
		return nil, true
	}

	body, err := doc.Html()
	if err != nil {
		log.Fatal(err)
	}

	enc, _, _ := charset.DetermineEncoding([]byte(body), "")
	decoder := enc.NewDecoder()
	x.Pages <- x.Parse(ctx, doc, decoder)
	return nil, true
}

func (x *CommentExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && x.Domain.MatchString(ctx.NormalizedURL().String())
}
