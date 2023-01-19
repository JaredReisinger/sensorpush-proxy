package jsonapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Wraps http.Client in an interface so that we can mock it for testing.
type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a JSON REST API client.
type Client struct {
	baseURL   string
	base      *url.URL
	requester HTTPRequestDoer
	// DefaultHeader http.Header
	header http.Header
}

// NewClient creates a new JSONAPIClient.
func NewClient(baseURL string) (*Client, error) {
	return NewClientWithDoer(baseURL, &http.Client{
		Timeout: 10 * time.Second,
	})
}

// NewClientWithDoer creates a new JSONAPIClient with a specific request doer.
func NewClientWithDoer(baseURL string, requestDoer HTTPRequestDoer) (*Client, error) {
	// sanity-check the baseURL?
	url, err := url.Parse(baseURL)
	if err != nil {
		// log.Print("???????")
		return nil, err
	}

	return &Client{
		baseURL:   baseURL,
		base:      url,
		requester: requestDoer,
		header:    make(http.Header),
	}, nil
}

func (c *Client) Header() http.Header {
	return c.header
}

// auto-unmarshal error bodies?

// Call makes a call!
func (c *Client) Call(relURL string, body interface{}, response interface{}) error {
	// calculate API endpoint...
	url, err := c.base.Parse(relURL)
	if err != nil {
		return err
	}
	// log.Printf("calculated API endpoint: %q", url)

	// TODO: handle no-body case (should this be a GET?)

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// log.Printf("\nposting to %q: %q", url, b)
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}

	// Add the default headers...
	req.Header = c.header.Clone()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.requester.Do(req)

	if err != nil {
		// log.Printf("WHAT???\n\n\t%+v", res)
		return err
	}
	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	// log.Printf("got JSON: %q", b)
	if err != nil {
		return err
	}

	// log.Printf("\n\n\tresponse: %+v\n\n\tbody: %q", res, b)

	// In any sane world (?!) we could rely entirely on the HTTP status code and
	// return (only) that when an error occurs... but there are many APIs,
	// SensorPush included, that return useful data in the body regardless of
	// the status code.  There are also cases where the API returns the wrong
	// status code, and the body is needed for a full understanding.  (For
	// instance SensorPush returns 400 with an "access denied" in the body
	// rather than just using 401.)
	//
	// The cleanest way to handle this is to unmarshal into the response object
	// *regardless* of the status code, so that it's available even if we end up
	// returning an error.  The caller can use an anonymous embedded struct for
	// the error fields.  (Alternatively, we can include the body bytes in the
	// returned error, for separate unmarshaling?)

	err = json.Unmarshal(b, response)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return &HTTPStatusError{StatusCode: res.StatusCode, Body: b}
	}

	return nil
}
