package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the structure that defines how the user should
// define the config.json file to make custom requests
type Config struct {
	Name     string
	Method   string
	URL      string
	Body     string
	Header   map[string][]string
	Host     string
	Form     string
	PostForm map[string][]string
}

func getConfig() (c []Config, err error) {
	raw, err := ioutil.ReadFile("./callit.json")
	if err != nil {
		return
	}
	return c, json.Unmarshal(raw, &c)
}
