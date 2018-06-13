package call

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
	raw, err := ioutil.ReadFile("callit.json")
	if err != nil {
		return
	}
	return c, json.Unmarshal(raw, &c)
}
