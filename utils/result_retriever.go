package utils
import (
	"net/http"
	"os"
	"fmt"
	"errors"
	"io/ioutil"
	"time"
)

// base URL for Googles Search REST API
const apiBaseURL string = "https://www.googleapis.com/customsearch/v1?"


// maps a google search result
// plus a unique id which identifies the result internally
// and the minified document
type Result struct{
	title string
	description string
	url string
	document string
	id int
}

// makes a google search for the query and writes minified documents
// of the 5 best results to the database.
// Returns
func processSearchQuery(query string, raw string) (map[int]Result, error){
	topResults, err := googleQuery(query)
	if err != nil {
		// TODO
	}
	c := make(chan map[string]string)
	go minifyResults(topResults, c)
	waitForAllMinifiedResults(c)
	return topResults, nil
}

// makes a query to the database for the requested result url.
// Returns empty string if nothing is found,
// otherwise a string containing the (minified) document
func processResultQuery(url string) (minifiedDoc string){
	requestedDocument := queryDBKey(url)
	return requestedDocument
}

// querys Googles REST search api.
// error is non-nill, if something failed.
// Else this returns a map containing the 5 best matches/results
// (in descending order).
func googleQuery(query string) (map[int]Result, error){
	apiKey := string(os.Args[2])
	cxID := string(os.Args[3])
	queryString := fmt.Sprint("key=", apiKey, "&cx=", cxID, "&q=", query)
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
	return
}

// waits till all 5 results are received on the channel
func waitForAllMinifiedResults(c chan map[string]string){
	returnedResults := 0
	for returnedResults < 5 {
		select {
		case minifiedResult := <-c:
			returnedResults++
			go writeToDB(minifiedResult)
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
	return
}

// writes result into redis DB to cache it.
// key: url, value: minified document
func writeToDB(map[string]string){
	//TODO
}

// queries one KEY of db, returns appropiate value if present
// or the empty string if not.
func queryDBKey(key string) (string){
	// TODO
	return
}
