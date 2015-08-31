package utils
import (
	"net/http"
	"gopkg.in/redis.v3"
)



type searchHandler struct {
	client *redis.Client
}

type resultHandler struct {
	client *redis.Client
}


// Starts a webserver which listens on
// '/' and '/results' route.
func StartWebserver(addr string) {
	redisClient := NewRedisClient()
	searchHandler := &searchHandler{client:redisClient}
	resultHandler := &resultHandler{client:redisClient}
	http.Handle("/", searchHandler)
	http.Handle("/results", resultHandler)
	http.ListenAndServe(addr, nil)
}

// Handles all search requests.
// These are not specified in the route
// but in the querystring.
// Parses query string, fails early if no
// 'search' parameter is provided.
func (sh searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if queryParams.Get("search") == "" {
		w.Write([]byte("Please specifiy a search in the query string"))
	}

	searchResult, err := processSearchQuery(queryParams.Get("search"),
		queryParams.Get("raw"), sh.client)
	if err != nil {
		// TODO
	}
	writeSearchResponse(w, searchResult)
}


// Handles all (subsequent) requests for search results.
// These are not specified in the route
// but in the querystring.
// Parses query string, fails early if no
// 'url' parameter is provided.
func (rh resultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if queryParams.Get("url") == "" {
		w.Write([]byte("No id specified."))
	}

	result := processResultQuery(queryParams.Get("id"), rh.client)
	if result == "" {
		// TODO
	}
	writeResultResponse(w, result)
}

// Writes the correspondent metadata results (url, title etc.)
// to the client, based on the first 5 results.
func writeSearchResponse(w http.ResponseWriter, result map[int]Result) {
	//	w.Header().Set("Content-Type", "text/html")
	//	w.Write([]byte("This is dog."))
	// TODO
}

// Writes the correspondent (minified) html
// for a single result to the client.
func writeResultResponse(w http.ResponseWriter, minHTML string) {
	// TODO
}