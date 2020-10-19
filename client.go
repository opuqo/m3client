package m3client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var defaultAPIURL = "api/v1/query"

type Client struct {
	url string
}

func NewClient(host string) *Client {
	return &Client{url: fmt.Sprintf("%s/%s", host, defaultAPIURL)}
}

func (c *Client) Request(param map[string]string) []byte {

	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}

	req.Header.Add("M3-Engine", "m3query")

	q := req.URL.Query()

	//http://localhost:7201/api/v1/query?query=count(http_requests)&time=1590147165

	for p, v := range param {
		q.Add(p, v)
	}

	req.URL.RawQuery = q.Encode()

	httpClient := &http.Client{}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	defer res.Body.Close()
	responce, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responce
}

func (c *Client) Count(tag string, timepoint time.Time) []byte {
	query := fmt.Sprintf("count(%s)", tag)

	param := make(map[string]string)
	param["query"] = query
	param["time"] = string(timepoint.Unix())

	return c.request(param)
}
func (c *Client) Max(tag string, timepoint time.Time) []byte {
	query := fmt.Sprintf("max(%s)", tag)

	param := make(map[string]string)
	param["query"] = query
	param["time"] = string(timepoint.Unix())

	return c.request(param)
}
