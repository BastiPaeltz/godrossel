package utils
import "gopkg.in/redis.v3"

// everything that has to do with displaying already retrieved
// documents/results goes in here


// makes a query to the database for the requested result url.
// Returns empty string if nothing is found,
// otherwise a string containing the (minified) document
func processResultQuery(url string, client *redis.Client) (minifiedDoc string){
	requestedDocument := queryDBKey(client, url)
	return requestedDocument
}