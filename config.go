package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the structure that defines how the user should
// define the config.json file to make custom requests
type Config struct {
	Name     string              `json:"name"`
	Method   string              `json:"method"`
	URL      string              `json:"url"`
	Body     string              `json:"body"`
	Header   map[string][]string `json:"header"`
	Host     string              `json:"host"`
	Form     string              `json:"form"`
	PostForm map[string][]string `json:"postform"`
}

func getConfig() (c []Config, err error) {
	raw, err := ioutil.ReadFile("./callit.json")
	if err != nil {
		return
	}
	return c, json.Unmarshal(raw, &c)
}
