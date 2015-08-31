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
)

// everything that has to do with getting results from the google search API
// and putting it into the redis DB and eventually returning the result to
// the caller goes in here.

// base URL for Googles Search REST API
const apiBaseURL string = "https://www.googleapis.com/customsearch/v1?"

// encapsulates a result from a google search
// and the minified document content (html content)
type Result struct{
	title string
	description string
	url string
	document string
}

// maps the returned JSON from the Google Search API
// (only interesting fields from response included)
// TODO
type SearchApiResponse struct{

}

// makes a google search for the query and writes minified documents
// of the 5 best results to the database.
// Returns
func processSearchQuery(query string, raw string, client *redis.Client) (map[int]Result, error){
	//TODO: process 'raw' parameter
	topResults, err := googleSearchQuery(query)
	if err != nil {
		// TODO
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
func googleSearchQuery(query string) (map[int]Result, error){
	apiKey := string(os.Args[2])
	cxID := string(os.Args[3])
	queryString := fmt.Sprint("key=", apiKey, "&cx=", cxID, "&q=", url.QueryEscape(query))
	apiResponse, err := http.Get(apiBaseURL + queryString)

	if err != nil || apiResponse.StatusCode != 200 {
		// TODO: write to log
		return nil, errors.New("Failed api response.")
	}

	responseBody, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		// TODO: write to log
		return nil, errors.New("Couldn't read response body.")
	}
	return filterResponse(responseBody), nil
}

// filters JSON response from the API
// returns a map containing the 5 best matches/results
// (in descending order)
func filterResponse(jsonBody []byte) (map[int]Result){
	return nil
}

// waits till all 5 results are received on the channel,
// writes to DB concurrently
func waitForAllMinifiedResults(client *redis.Client, c chan map[string]string){
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

