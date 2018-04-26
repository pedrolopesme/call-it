package call

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
