package utils
import (
	"github.com/PuerkitoBio/goquery"
	"log"
)


// Minifies search result documents by making a GET request to the URL,
// if it is not cached already. Writes result to channel.
func minifyResults(resMap *[]Result, c chan map[string]string) {
	for _, result := range *resMap {
		go minifyDocument(result.Link, c)
	}
}

// Minifies document, means no media content, no css, no js
func minifyDocument(url string, c chan map[string]string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		c <- map[string]string{url: "Site not available"}
		return
	}
	tagsToBeRemoved := []string{"img", "video", "link", "audio", "track",
		"embed", "iframe", "source", "canvas", "script"}
	for _, singleTag := range tagsToBeRemoved {
		doc.Find(singleTag).Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
	}
	minifiedHtml, err := doc.Html()
	if err != nil {
		log.Println("getting html from goquery doc failed.", err.Error())
		c <- map[string]string{url: "Site not available"}
	}
	c <- map[string]string{url: minifiedHtml}
}


