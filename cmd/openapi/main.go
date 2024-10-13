package main

import (
	"embed"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed proto/gen/openapiv3/*
var fs embed.FS

func main() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, uri *url.URL) ([]byte, error) {
		return fs.ReadFile(uri.Path)
	}

	doc, err := loader.LoadFromFile("proto/gen/openapiv3/openapi.yml")
	if err != nil {
		panic(err)
	}

	if err = doc.Validate(loader.Context); err != nil {
		panic(err)
	}
}
