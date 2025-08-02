package call

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMakeCallsWhenURLExists(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 100)
	result := call.MakeIt()

	assert.Equal(test, 1, len(result.status))
	assert.Equal(test, 10, result.status[200].total)
}

func TestMakeCallsWhenURLDoesntExist(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(404, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 100)
	result := call.MakeIt()

	assert.Equal(test, 1, len(result.status))
	assert.Equal(test, 0, result.status[200].total)
	assert.Equal(test, 10, result.status[404].total)
}

func TestMakeCallsReturnTheSameStatusCode(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "100"}
	call, _ := BuildCall(params, 1, 10)
	result := call.MakeIt()

	assert.Equal(test, 1, len(result.status))
	assert.Equal(test, 100, result.status[200].total)
}

func TestCalcConcurrentAttemptsWhenThereAreEnoughAttemptsLeft(test *testing.T) {
	urlAddress, _ := url.Parse("http://www.a.com")
	call := ConcurrentCall{URL: urlAddress, Attempts: 100, ConcurrentAttempts: 10}

	assert.Equal(test, 10, calcConcurrentAttempts(call))
}

func TestCalcConcurrentAttemptsWhenThereAreNotEnoughAttemptsLeft(test *testing.T) {
	urlAddress, _ := url.Parse("http://www.a.com")
	call := ConcurrentCall{URL: urlAddress, Attempts: 10, ConcurrentAttempts: 100}

	assert.Equal(test, 10, calcConcurrentAttempts(call))
}

func TestCalcConcurrentAttemptsWhenAttemptsLeftIsEqualToConcurrentAttempts(test *testing.T) {
	urlAddress, _ := url.Parse("http://www.a.com")
	call := ConcurrentCall{URL: urlAddress, Attempts: 10, ConcurrentAttempts: 10}

	assert.Equal(test, 10, calcConcurrentAttempts(call))
}

func TestGetUrl(test *testing.T) {
	urlAddress := "http://www.foo.com/bar"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", urlAddress,
		httpmock.NewStringResponder(200, `[]`))

	config := Config{}
	parsedURL, _ := url.Parse(urlAddress)
	callResponses := callURL(parsedURL, 50, config)

	for _, response := range callResponses {
		assert.Equal(test, 200, response.status)
	}
	assert.Equal(test, 50, len(callResponses))
}

func Test_buildRequest(t *testing.T) {
	type args struct {
		baseURL string
		config  Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should pass without config",
			args: args{
				baseURL: "https://www.survivingmars.com",
			},
			wantErr: false,
		},
		{
			name: "should pass with config",
			args: args{
				baseURL: "https://www.survivingmars.com",
				config: Config{
					Name:   "elon musk",
					Method: http.MethodPost,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildRequest(tt.args.baseURL, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
