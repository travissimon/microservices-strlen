package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/travissimon/remnant/client"
)

func swaggerHandler(w http.ResponseWriter, r *http.Request, logger client.Logger) {
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

const (
	SERVICE_PATH = "/v1/strlen/"
)

func strLenHandler(w http.ResponseWriter, r *http.Request, logger client.Logger) {
	logger.LogDebug("request: " + r.URL.Path)
	s := string(r.URL.Path[len(SERVICE_PATH):])
	len := utf8.RuneCountInString(s)
	logger.LogInfo(fmt.Sprintf("Str len: %d", len))
	fmt.Fprintf(w, "%d", len)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("healthz %s\n", time.Now().Local())
	fmt.Fprintf(w, "OK")
}

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")
	flag.Parse()

	remnantUrl := "http://localhost:7777/"
	http.HandleFunc(SERVICE_PATH, client.GetInstrumentedHandler(remnantUrl, strLenHandler))
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/swagger.json", client.GetInstrumentedHandler(remnantUrl, swaggerHandler))
	http.HandleFunc("/", client.GetInstrumentedHandler(remnantUrl, swaggerHandler))

	fmt.Printf("Starting server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}
