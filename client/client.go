package client

import (
	"io"
	"net/http"
)

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{}
}

func (c *Client) Fetch(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "drafter/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
