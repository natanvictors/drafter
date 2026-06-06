package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/natanvictors/drafter/parser"
)

type Client struct {
	http *http.Client
}

func New() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &Client{
		http: &http.Client{Jar: jar},
	}, nil
}

func (c *Client) Fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "drafter/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) Login() error {

	resp, err := c.http.Get("https://lol.fandom.com/api.php?action=query&meta=tokens&type=login&format=json")
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := parser.ParseToken(body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	resp, err = c.http.PostForm("https://lol.fandom.com/api.php", url.Values{
		"action":     {"login"},
		"lgname":     {os.Getenv("API_BOT_NAME")},
		"lgpassword": {os.Getenv("API_BOT_PASSWORD")},
		"lgtoken":    {data.Query.Tokens.LoginToken},
		"format":     {"json"},
	})
	if err != nil {
		return err
	}

	return nil
}
