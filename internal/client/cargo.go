package cargo

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/natanvictors/drafter/internal/parser"
)

type Client struct {
	http     *http.Client
	user     string
	password string
}

type RequestInfo struct {
	URL        string
	Maxretries int
	retries    int
	Interval   int
}

func (info *RequestInfo) addRetry() *RequestInfo {
	info.retries++
	return info
}

func New(user string, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &Client{
		http:     &http.Client{Jar: jar},
		user:     user,
		password: password,
	}, nil
}

func (c *Client) FetchAndParse(info RequestInfo) (parser.Cargoresult, error) {
	req, err := http.NewRequest("GET", info.URL, nil)
	if err != nil {
		return parser.Cargoresult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "drafter/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return parser.Cargoresult{}, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	result, err := parser.ParseCargo(data)
	if err != nil {
		return parser.Cargoresult{}, fmt.Errorf("parse failed: %s", err)
	}

	if result.Cargoerror != nil {
		if info.retries > info.Maxretries {
			return parser.Cargoresult{}, fmt.Errorf("couldnt fetch after %d retries", info.Maxretries)
		}
		fmt.Println("fetch failed:", result.Cargoerror.Code, ". retrying...")
		info.addRetry()
		time.Sleep(time.Duration(info.Interval) * time.Second)
		return c.FetchAndParse(info)
	}

	return *result, nil
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

	_, err = c.http.PostForm("https://lol.fandom.com/api.php", url.Values{
		"action":     {"login"},
		"lgname":     {c.user},
		"lgpassword": {c.password},
		"lgtoken":    {data.Query.Tokens.LoginToken},
		"format":     {"json"},
	})
	if err != nil {
		return err
	}

	return nil
}
