package crawler

import (
	"log"
	"time"
)

func waitForProcessingOfAllPages(pages <-chan Page, processPage func(Page)) {
	for {
		select {
		case page := <-pages:
			processPage(page)
		case <-time.After(1 * time.Second):
			log.Println("Processed all pages")
			return
		}
	}
}

func WaitForPages(pages <-chan Page, doneCrawling <- chan bool, processPage func(Page)) {
	for {
		select {
		case page := <-pages:
			processPage(page)
		case <-doneCrawling:
			waitForProcessingOfAllPages(pages, processPage)
			return
		}
	}
}
