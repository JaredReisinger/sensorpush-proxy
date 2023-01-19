package sensorpush

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/jsonapi"
)

type ApiClient interface {
	Header() http.Header
	Call(url string, body interface{}, response interface{}) error
}

// Client is a SensorPush client
type Client struct {
	email       string
	password    string
	authToken   string
	accessToken string

	// client *jsonapi.Client
	client ApiClient
}

const apiBase = "https://api.sensorpush.com/api/v1/"

// NewClient creates a new SensorPush client
func NewClient(email string, password string) (*Client, error) {
	client, err := jsonapi.NewClient(apiBase)
	if err != nil {
		return nil, err
	}

	return NewClientWithApiClient(email, password, client)
}

// NewClientWithApiClient creates a new SensorPush client
func NewClientWithApiClient(email string, password string, client ApiClient) (*Client, error) {
	return &Client{
		email:    email,
		password: password,
		client:   client,
	}, nil
}

// SensorPush doesn't understand how to properly use HTTP status codes... they
// return 400 (Bad Request) for ACCESS_DENIED, when they should be using 401
// (Unauthorized).  We *could* try to extract that information, but meh?

type errorResponse struct {
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode string `json:"statusCode"`
}

type authorizeRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authorizeResponse struct {
	Authorization string `json:"authorization"`
	APIkey        string `json:"apikey"`
}

type accessTokenRequest struct {
	Authorization string `json:"authorization"`
}

type accessTokenResponse struct {
	AccessToken string `json:"accesstoken"`
}

func (c *Client) authorize() error {
	c.client.Header().Del("Authorization")

	authorize := &authorizeResponse{}
	err := c.client.Call("oauth/authorize", &authorizeRequest{Email: c.email, Password: c.password}, authorize)
	if err != nil {
		log.Printf("unable to authorize: %+v", err)
		return err
	}
	// log.Printf("got auth: %+v", authorize)
	c.authToken = authorize.Authorization
	return nil
}

func (c *Client) access(depth int) error {
	var err error

	if c.authToken == "" {
		err = c.authorize()
		if err != nil {
			return err
		}
	}

	c.client.Header().Del("Authorization")

	accessToken := &accessTokenResponse{}
	err = c.client.Call("oauth/accesstoken", &accessTokenRequest{Authorization: c.authToken}, accessToken)
	// exit early on success
	if err == nil {
		// log.Printf("got access token: %+v", accessToken)
		c.accessToken = accessToken.AccessToken
		// also update the default headers...
		c.client.Header().Set("Authorization", c.accessToken)
		return nil
	}

	log.Printf("unable to get access token: %+v", err)
	// TODO: check if authorization token expired and we need to sign in
	// again?
	// now call recursively unless we've reached our maximum depth!
	if depth < 1 {
		return c.access(depth + 1)
	}

	return err
}

func (c *Client) authCall(relURL string, input interface{}, output interface{}, depth int) error {
	var err error

	if c.accessToken == "" {
		err = c.access(0)
		if err != nil {
			return err
		}
	}

	err = c.client.Call(relURL, input, output)
	// exit early on success
	if err == nil {
		return nil
	}

	switch t := err.(type) {
	case *jsonapi.HTTPStatusError:
		errorBody := &errorResponse{}
		json.Unmarshal(t.Body, errorBody) // ignore error?
		log.Printf("alternate view of error (%s): %#v", t.Error(), errorBody)
		if errorBody.Type == "ACCESS_DENIED" {
			c.accessToken = ""
			// now call recursively unless we've reached our maximum depth!
			if depth < 2 {
				return c.authCall(relURL, input, output, depth+1)
			}
		}
	default:
		log.Printf("unable to get %q: %+v", relURL, err)
	}
	return err
}
