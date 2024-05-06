package callemall

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mrjones/oauth"
)

// Client is our type
type Client struct {
	endpoint    string
	httpclient  *http.Client
	prevRequest time.Time
}

// NewClient creates a new call-em-all client
func NewClient(appkey, secretkey, authtoken, endpoint string) (*Client, error) {
	// setup for oauth signing
	consumer := oauth.NewConsumer(appkey, secretkey, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: authtoken}

	// get http client that handles oauth
	httpclient, err := consumer.MakeHttpClient(accessToken)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	client := &Client{
		endpoint:   endpoint,
		httpclient: httpclient,
	}

	return client, nil
}

// makeRequest is a helper function to wrap making REST calls to call-em-all
func (client *Client) makeRequest(method, url string, body io.Reader) ([]byte, error) {
	// call-em-all does rate limiting
	// make sure it has been at least 100 milliseconds since the last call
	elapse := time.Since(client.prevRequest)
	if elapse < 100*time.Millisecond {
		time.Sleep(100 * time.Millisecond)
	}
	defer func() { client.prevRequest = time.Now() }()

	// create request
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	// set content-type only on requests that send some content
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	// make request, get response
	var response *http.Response
	response, err = client.httpclient.Do(request)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	defer response.Body.Close()

	// error?
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s call to %s returned status code %d ", method, url, response.StatusCode)
		log.Printf("%+v", err)
		return nil, err
	}

	// get body for caller, if there is something
	var data []byte
	if response.ContentLength != 0 {
		data, err = io.ReadAll(response.Body)
		if err != nil {
			log.Printf("%+v", err)
			return nil, err
		}
	}

	return data, nil
}
