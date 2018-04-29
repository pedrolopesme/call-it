package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
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
	dir, err := os.Open(".")
	if err != nil {
		return
	}
	files, err := dir.Readdir(0)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.Name() != "callit.json" {
			continue
		}
		raw, errRaw := ioutil.ReadFile(f.Name())
		if err != nil {
			return c, errRaw
		}
		return c, json.Unmarshal(raw, &c)
	}
	return c, errors.New("Could not find file callit.json")
}
