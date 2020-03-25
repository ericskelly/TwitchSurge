package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// GetJSON - For making get requests
func GetJSON(url string, output interface{}) error {
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(output)
}

// ReadJSONConfig - read the json config file
func ReadJSONConfig(configuration *Config) error {
	file, err := os.Open("../../conf.json")
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(configuration)
}
