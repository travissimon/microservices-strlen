package main

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	respJson := `{
    "kind": "pathList",
    "items": [
        "/v1",
        "/v2",
        "/healthz"
    ],
    "documentation": "http://github.com/travissimon/microservices/strlen/README.md",
    "contact": "travis.simon@nicta.com.au"
}`

	fmt.Fprintf(w, respJson)
}

func defaultStrLenHandler(w http.ResponseWriter, r *http.Request) {
	respJson := `{
    "kind": "pathList",
    "items": [
        "/strlen",
    ],
}`

	fmt.Fprintf(w, respJson)
}

func strLenV1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: %s\n", r.URL.Path)
	s := string(r.URL.Path[len("/v1/len/"):])
	fmt.Fprintf(w, "%d", utf8.RuneCountInString(s))
}

func strLenV2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: %s\n", r.URL.Path)
	s := string(r.URL.Path[len("/v2/len/"):])
	fmt.Fprintf(w, "%d", utf8.RuneCountInString(s))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("healthz %s\n", time.Now().Local())
	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/v1/", defaultStrLenHandler)
	http.HandleFunc("/v2/", defaultStrLenHandler)
	http.HandleFunc("/v1/len/", strLenV1Handler)
	http.HandleFunc("/v2/len/", strLenV2Handler)
	http.HandleFunc("/healthz", healthzHandler)

	fmt.Printf("Starting server on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
