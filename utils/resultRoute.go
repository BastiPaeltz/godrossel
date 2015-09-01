package utils
import (
	"gopkg.in/redis.v3"
	"errors"
)

// everything that has to do with displaying already retrieved
// documents/results goes in here (requesting '/result')


// makes a query to the database for the requested result url.
// Returns empty string if nothing is found,
// otherwise a string containing the (minified) document
func processResultQuery(url string, client *redis.Client) (minifiedDoc string, err error){
	requestedDocument := queryDBKey(client, url)
	if requestedDocument == ""{
		return "", errors.New("No document entry for this URL.")
	}
	return requestedDocument, nil
}