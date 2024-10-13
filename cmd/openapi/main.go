package main

import (
	"embed"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed proto/gen/openapiv3/* swagger-ui/*
var fs embed.FS

func main() {
	// Validate the spec first
	validateSpec()

	// Serve Swagger UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Serve OpenAPI files directly
		if strings.HasPrefix(path, "/api/") {
			content, err := fs.ReadFile(filepath.Join("proto/gen/openapiv3", path[5:]))
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				w.Header().Set("Content-Type", "application/yaml")
			} else {
				w.Header().Set("Content-Type", "application/json")
			}
			w.Write(content)
			return
		}

		// Serve Swagger UI files
		content, err := fs.ReadFile(filepath.Join("swagger-ui", path))
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		w.Write(content)
	})

	log.Println("Documentation server started at http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func validateSpec() {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("proto/gen/openapiv3/openapi.yaml")
	if err != nil {
		panic(err)
	}

	if err = doc.Validate(loader.Context); err != nil {
		panic(err)
	}
}
