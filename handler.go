// Package urlshort is  for url handler functions
package urlshort

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"
)

func MySQLHandler(fallback http.Handler) http.HandlerFunc {
	// connecting to the MySQL DB
	db, err := sql.Open("mysql", "root:newpassword@tcp(localhost:3306)/url_paths")
	if err != nil {
		fmt.Println("error validating sql.Open arguments")
		panic(err)
	}
	// defer db.Close()

	// ping DB to check connection
	err = db.Ping()
	if err != nil {
		fmt.Println("error veryfing connection with db.Ping")
		panic(err)
	}
	fmt.Println("Succesfull connection to the database!")

	return func(w http.ResponseWriter, r *http.Request) {
		var url string
		path := strings.Trim(r.URL.Path, "/")
		err = db.QueryRow("SELECT url FROM paths where path = ?", path).Scan(&url)

		if err == sql.ErrNoRows {
			fallback.ServeHTTP(w, r)
			return
		} else if err != nil {
			panic(err)
		} else {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
	}
}

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
