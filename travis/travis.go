package travis

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion          = "0.1"
	defaultBaseURL          = "https://api.travis-ci.org/"
	userAgent               = "go-travis/" + libraryVersion
	defaultTravisAPIVersion = "3"
)

type Client struct {
	client            *http.Client
	BaseURL           *url.URL
	UserAgent         string
	TravisAPIVersion  string
	Repository        *RepositoryService
	TravisAccessToken string
}

func NewClient(travisAccessToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, TravisAccessToken: travisAccessToken}

	//TODO add all other services
	c.Repository = &RepositoryService{client: c}

	return c
}

func (c *Client) NewRequest(method, urlStr string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		println("Error in url parse")
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		println("Error in http new request")
		return nil, err
	}

	req.Header.Add("Travis-API-Version", c.TravisAPIVersion)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", "token "+c.TravisAccessToken)
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		println("Error in http client do")
		return nil, err
	}
	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()
	response := resp //TODO: needs own response with pagination

	err = CheckResponse(resp)
	if err != nil {
		println("Check Response failed")
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Cannot read resp.Body")
		return response, err
	}
	json.Unmarshal(bytes, &v)
	return response, err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		println("HTTP Status 2XX")
		return nil
	}
	println("HTTP Status not in 2XX")
	return errors.New("HTTP Response: " + strconv.Itoa(r.StatusCode))
}
