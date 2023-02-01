package pkg

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	comingSoonAction = "Prochainement"
	reserveAction    = "Réserver"
	monacoPSGMatch   = "AS MONACOPARIS SAINT GERMAIN"
)

func Scrape(notify func()) {
	// Request the HTML page.
	res, err := http.Get(MonacoBilleterieURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// "exists(.matchCard|text(span.teamName) == 'AS MONACOPARIS SAINT GERMAIN' && text(.matchActions) == 'Réserver')"

	spaceReducer := regexp.MustCompile(`\s+`)
	newLineReducer := regexp.MustCompile(`\n+`)

	// Find the available sell options
	doc.Find(".matchPrdList").Each(func(i int, s *goquery.Selection) {
		// count the acheter buttons
		cleanedFullText := newLineReducer.ReplaceAllString(spaceReducer.ReplaceAllString(s.Text(), " "), " ")
		log.Printf("full text: %s", cleanedFullText)
		if strings.Contains(cleanedFullText, "Grand Public Réserver") || !strings.Contains(cleanedFullText, "Grand Public A partir du 03/02/2023") {
			notify()
		}
	})
}
