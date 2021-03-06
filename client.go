package redmine

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	endpoint string
	apikey   string
	Limit    int
	Offset   int
	*http.Client
}

const NoSetting = -1

var DefaultLimit int = NoSetting
var DefaultOffset int = NoSetting

func NewClient(endpoint, apikey string) *Client {
	return &Client{
		endpoint: endpoint,
		apikey:   apikey,
		Limit:    DefaultLimit,
		Offset:   DefaultOffset,
		Client:   http.DefaultClient,
	}
}

func (c *Client) apiKeyParameter() string {
	return "key=" + c.apikey
}

func (c *Client) concatParameters(requestParameters ...string) string {
	cleanedParams := []string{}
	for _, param := range requestParameters {
		if param != "" {
			cleanedParams = append(cleanedParams, param)
		}
	}

	return strings.Join(cleanedParams, "&")
}

// URLWithFilter return string url by concat endpoint, path and filter
// err != nil when endpoint can not parse
func (c *Client) URLWithFilter(path string, f Filter) (string, error) {
	var fullURL *url.URL
	fullURL, err := url.Parse(c.endpoint)
	if err != nil {
		return "", err
	}
	fullURL.Path += path
	if c.Limit > -1 {
		f.AddPair("limit", strconv.Itoa(c.Limit))
	}
	if c.Offset > -1 {
		f.AddPair("offset", strconv.Itoa(c.Offset))
	}
	fullURL.RawQuery = f.ToURLParams()
	return fullURL.String(), nil
}

func (c *Client) getPaginationClause() string {
	clause := ""
	if c.Limit > -1 {
		clause = clause + fmt.Sprintf("&limit=%d", c.Limit)
	}
	if c.Offset > -1 {
		clause = clause + fmt.Sprintf("&offset=%d", c.Offset)
	}
	return clause
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

type IdName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Id struct {
	Id int `json:"id"`
}
