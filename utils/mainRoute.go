package utils
import (
	"net/http"
	"os"
	"fmt"
	"errors"
	"io/ioutil"
	"time"
	"net/url"
	"gopkg.in/redis.v3"
	"log"
	"encoding/json"
)

// everything that has to do with getting results from the google search API
// and putting it into the redis DB and eventually returning the result to
// the caller goes in here.
// (requesting '/')

// base URL for Googles Search REST API
const apiBaseURL string = "https://www.googleapis.com/customsearch/v1?"

// encapsulates a result from a google search
// and the minified document content (html content)
type Result struct {
	Title       string
	HtmlTitle   string
	Link        string
	DisplayLink string
	Snippet     string
	HtmlSnippet string
}

// maps the returned JSON from the Google Search API
// (only interesting fields from response included)
type SearchApiResponse struct {
	Items []Result
}


// makes a google search for the query and writes minified documents
// of the 5 best results to the database.
// Returns
func processSearchQuery(query string, raw string, client *redis.Client) (*[]Result, error) {
	//TODO: process 'raw' parameter
	topResults, err := googleSearchQuery(query)
	if err != nil {
		log.Println("search query failed.")
	}
	c := make(chan map[string]string)
	go minifyResults(topResults, c)
	waitForAllMinifiedResults(client, c)
	return topResults, nil
}

// querys Googles REST search api.
// error is non-nil, if something failed.
// Else this returns a map containing the 5 best matches/results
// (in descending order).
func googleSearchQuery(query string) (*[]Result, error) {
	cxID, _ := os.LookupEnv("CXID")
	apiKey, _ := os.LookupEnv("APIKEY")

	queryString := fmt.Sprint("key=", apiKey, "&cx=", cxID, "&q=", url.QueryEscape(query), "&num=5")
	apiResponse, err := http.Get(apiBaseURL + queryString)

	if err != nil || apiResponse.StatusCode != 200 {
		log.Println("Failed api response (status code != 200).")
		return nil, errors.New("Failed api response.")
	}

	responseBody, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		log.Println("Couldn't read api response.")
		return nil, errors.New("Couldn't read response body.")
	}
	return filterResponse(responseBody), nil
}

// filters JSON response from the API
// returns a map containing the 5 best matches/results
// (in descending order)
func filterResponse(jsonBody []byte) (*[]Result) {
	var resp SearchApiResponse
	err := json.Unmarshal(jsonBody, &resp)
	if err != nil {
		log.Println("Cant parse/unmarshal json response.")
	}
	return &resp.Items
}

// waits till all 5 results are received on the channel,
// writes to DB concurrently
func waitForAllMinifiedResults(client *redis.Client, c chan map[string]string) {
	returnedResults := 0
	for returnedResults < 5 {
		select {
		case minifiedResult := <-c:
			returnedResults++
			go writeToDB(client, minifiedResult)
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
	return
}

