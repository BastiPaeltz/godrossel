package utils
import (
	"net/http"
	"gopkg.in/redis.v3"
	"text/template"
	"log"
)

var renderedTemplate = template.Must(template.ParseFiles("results.html", "notFound.html"))

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
		writeNotFoundResponse(w)
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
	if queryParams.Get("url") == "" {
		writeNotFoundResponse(w)
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
		log.Println("Template rendering error for results.html", err.Error())
	}
}

// Writes the correspondent (minified) html
// for a single result to the client.
func writeResultResponse(w http.ResponseWriter, minHTML string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(minHTML))
}

// Writes html page with a "Resource not found" error
// back to the client.
func writeNotFoundResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	err := renderedTemplate.ExecuteTemplate(w, "notFound.html", nil)
	if err != nil {
		log.Println("Template rendering error for notFound.html", err.Error())
	}
}