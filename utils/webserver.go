package utils
import "net/http"


// Starts a webserver which only listens to
// '/' and '/results' route.
func StartWebserver(addr string) {
	http.HandleFunc("/", search)
	http.HandleFunc("/results", results)
	http.ListenAndServe(addr, nil)
}

// Handles all search requests.
// These are not specified in the route
// but in the querystring.
// Parses query string, fails early if no
// 'search' parameter is provided.
func search (w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query()
	if queryParams.Get("search") == ""{
		w.Write([]byte("Please specifiy a search in the query string"))
	}

	searchResult, err := processSearchQuery(queryParams.Get("search"),
		queryParams.Get("raw"))
	if err != nil{
		// TODO
	}
	writeSearchResponse(&w, searchResult)
}


// Handles all (subsequent) requests for search results.
// These are not specified in the route
// but in the querystring.
// Parses query string, fails early if no
// 'url' parameter is provided.
func results (w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query()
	if queryParams.Get("url") == ""{
		w.Write([]byte("No id specified."))
	}

	result := processResultQuery(queryParams.Get("id"))
	if result == "" {
		// TODO
	}
	writeResultResponse(&w, result)
}

// Writes the correspondent metadata results (url, title etc.)
// to the client, based on the first 5 results.
func writeSearchResponse (w *http.ResponseWriter, result map[int]Result ){

}

// Writes the correspondent (minified) html
// for a single result to the client.
func writeResultResponse (w *http.ResponseWriter, minHTML string ){

}