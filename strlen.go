package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/travissimon/remnant/client"
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	respJson := `{
    "swagger": "2.0",
    "info": {
        "description": "This finds the length of a string. It's AWESOME!\n",
        "version": "1.0.0",
        "title": "strlen",
        "contact": {
            "name": "Travis"
            "email": "travis.simon@nicta.com.au"
        },
        "license": {
            "name": "GPL 3.0",
            "url": "http://www.gnu.org/licenses/gpl-3.0.en.html"
        }
    },
    "host": "strlen",
    "basePath": "/v1",
    "schemes": [
        "http"
    ],
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "paths": {
        "/strlen/{string}": {
            "get": {
                "tags": [
                    "string"
                ],
                "summary": "Find the length of a string",
                "description": "Returns length of a string",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "string",
                        "description": "string whose length is unknown",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successful operation",
                        "schema": {
                            "type": "integer",
                            "format": "int64"
                        }
                    }
                }
            }
        },
        "/healthz": {
            "get": {
                "tags": [
                    "monitoring"
                ],
                "summary": "Heartbeat response",
                "description": "Provides a monitoring service with proof-of-life",
                "operationId": "isAlive",
                "responses": {
                    "default": {
                        "description": "successful operation"
                    }
                }
            }
        }
    }
}`

	fmt.Fprintf(w, respJson)
}

func strLenHandler(w http.ResponseWriter, r *http.Request) {
	cl, err := client.NewRemnantClient("http://localhost:7777/", r)
	defer cl.EndSpan()

	fmt.Printf("%s request: %s\n", time.Now().UTC().Format(time.RFC3339), r.URL.Path)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	s := string(r.URL.Path[len("/v1/strlen/"):])
	fmt.Fprintf(w, "%d", utf8.RuneCountInString(s))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("healthz %s\n", time.Now().Local())
	fmt.Fprintf(w, "OK")
}

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")
	flag.Parse()

	http.HandleFunc("/", swaggerHandler)
	http.HandleFunc("/swagger.json", swaggerHandler)
	http.HandleFunc("/v1/strlen/", strLenHandler)
	http.HandleFunc("/healthz", healthzHandler)

	fmt.Printf("Starting server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}
