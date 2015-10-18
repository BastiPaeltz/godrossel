package utils
import (
	"net/http"
	"gopkg.in/redis.v3"
	"text/template"
	"log"
)

var renderedTemplate = template.Must(template.ParseFiles("results.html"))

type templateData struct {
	Result []Result
	Search string
}

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
	searchTerm := queryParams.Get("search")
	if searchTerm == "" {
		w.Write([]byte("<h1>Please specifiy a search in the query string</h1>"))
	} else {
		searchResult, err := processSearchQuery(searchTerm,
			queryParams.Get("raw"), sh.client)
		if err != nil {
			log.Println("Failed to process search.")
		}
		writeSearchResponse(w, *searchResult, searchTerm)
	}
}


// Handles all (subsequent) requests for search results.
// These are not specified in the route
// but in the querystring.
// Parses query string, fails early if no
// 'url' parameter is provided.
func (rh resultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	w.Header().Set("Content-Type", "text/html")
	if queryParams.Get("url") == "" {
		w.Write([]byte("<h1>No id specified!</h1>"))
	} else {
		result, err := processResultQuery(queryParams.Get("url"), rh.client)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		writeResultResponse(w, result)
	}
}

// Writes the correspondent metadata results (url, title etc.)
// to the client, based on the first 5 results.
func writeSearchResponse(w http.ResponseWriter, results []Result, searchTerm string) {
	w.Header().Set("Content-Type", "text/html")
	err := renderedTemplate.ExecuteTemplate(w, "results.html", &templateData{results, searchTerm})
	if err != nil {
		log.Println("Template rendering error", err.Error())
	}
}

// Writes the correspondent (minified) html
// for a single result to the client.
func writeResultResponse(w http.ResponseWriter, minHTML string) {
	w.Write([]byte(minHTML))
}