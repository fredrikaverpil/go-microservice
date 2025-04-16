package server

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SwaggerHandler creates an HTTP handler that serves Swagger UI and the OpenAPI specification.
func SwaggerHandler(logger *slog.Logger) http.Handler {
	// Define paths
	uiDir := "./swagger-ui"
	specPath := "./proto/gen/openapiv3/openapi.yaml"

	// Verify paths exist
	if _, err := os.Stat(uiDir); os.IsNotExist(err) {
		logger.Error("Swagger UI directory not found", "path", uiDir, "error", err)
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "Swagger UI files not found. Run 'make openapi-download-swagger-ui-files'",
				http.StatusInternalServerError)
		})
	}

	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		logger.Error("OpenAPI spec file not found", "path", specPath, "error", err)
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "OpenAPI specification file not found", http.StatusInternalServerError)
		})
	}

	fileServer := http.FileServer(http.Dir(uiDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the OpenAPI spec file
		if r.URL.Path == "/api/openapi.yaml" {
			w.Header().Set("Content-Type", "application/yaml")
			http.ServeFile(w, r, specPath)
			return
		}

		// Handle /docs and /docs/ by directly serving index.html
		if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
			http.ServeFile(w, r, filepath.Join(uiDir, "index.html"))
			return
		}

		// For other /docs/* paths, strip the prefix and let the file server handle it
		if strings.HasPrefix(r.URL.Path, "/docs/") {
			// This handles JS, CSS, and other static files
			http.StripPrefix("/docs/", fileServer).ServeHTTP(w, r)
			return
		}

		// Return 404 for any other path
		http.NotFound(w, r)
	})
}
