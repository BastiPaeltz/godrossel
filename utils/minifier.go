package utils
import (
	"net/http"
	"io/ioutil"
)

// Minifies search result documents by making a GET request to the URL,
// if it is not cached already. Writes result to channel.
func minifyResults(resMap map [int]Result,c chan map[string]string){
	for _, result := range resMap {
		if siteIsAlreadyCached(result.url){
			c <- map[string]string{"": ""}
			continue
		}

		resp, err := http.Get(result.url)
		if err != nil {
			// TODO
		}
		document, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// TODO
		}
		go minifyDocument(result.url, string(document), c)
	}
}

// Minifies document, means no media content, no css, no js
func minifyDocument(url string,document string, c chan map[string]string){
	var minifiedDoc string
	// TODO : minify document with goquery
	c <- map[string]string{url: minifiedDoc}
}

// returns true if site is already present in database,
// false otherwise
func siteIsAlreadyCached(url string) (bool){
	if queryDBKey(url) == ""{
		return false
	}
	return true
}