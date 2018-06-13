package call

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// Config is the structure that defines how the user should
// define the config.json file to make custom requests
type Config struct {
	Name               string              `json:"name,omitempty"`
	Method             string              `json:"method"`
	Attempts           int                 `json:"attempts,omitempty"`
	ConcurrentAttempts int                 `json:"concurrent,omitempty"`
	URL                string              `json:"url"`
	Body               string              `json:"body,omitempty"`
	Header             map[string][]string `json:"header,omitempty"`
	Host               string              `json:"host"`
	Form               string              `json:"form,omitempty"`
	PostForm           map[string][]string `json:"postform,omitempty"`
}

func (c *Config) checkDefaults() (err error) {
	if len(c.URL) == 0 {
		return errors.New("Cannot create call to empty url")
	}
	if c.Attempts == 0 {
		c.Attempts = 10
	}
	if c.ConcurrentAttempts == 0 {
		c.ConcurrentAttempts = 10
	}
	return
}

func config() (c []Config, err error) {
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
