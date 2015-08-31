package utils
import (
	"github.com/PuerkitoBio/goquery"
)


// Minifies search result documents by making a GET request to the URL,
// if it is not cached already. Writes result to channel.
func minifyResults(resMap map [int]Result,c chan map[string]string){
	for _, result := range resMap {
		doc, err := goquery.NewDocument(result.url)
		if err != nil {
			// TODO
		}
		go minifyDocument(result.url, doc, c)
	}
}

// Minifies document, means no media content, no css, no js
func minifyDocument(url string,doc *goquery.Document, c chan map[string]string){
	tagsToBeRemoved := []string{"img", "video", "link", "audio", "track",
		"embed", "iframe", "source", "canvas", "script"}
	for _, singleTag := range tagsToBeRemoved {
		doc.Find(singleTag).Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
	}
	minifiedHtml, err := doc.Html()
	if err != nil{
		//TODO
	}
	c <- map[string]string{url: minifiedHtml}
}


