// Package urlshort is  for url handler functions
package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if value, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, value, http.StatusFound)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func buildMap(mappings []URLMapping) map[string]string {
	pathToURLMap := make(map[string]string)

	for _, mapping := range mappings {
		pathToURLMap[mapping.Path] = mapping.URL
	}

	return pathToURLMap
}

type URLMapping struct {
	Path string `json:"path" yaml:"path"`
	URL  string `json:"url" yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYaml(yml []byte) ([]URLMapping, error) {
	var mappings []URLMapping
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		panic(err)
	}

	return mappings, nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsonFile []byte) ([]URLMapping, error) {
	var mappings []URLMapping
	err := json.Unmarshal(jsonFile, &mappings)
	if err != nil {
		panic(err)
	}

	return mappings, nil
}
