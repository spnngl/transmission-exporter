// Package transmission is a client of the Transmission RPC API, limited to what the exporter reads.
package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	endpoint        = "/transmission/rpc/"
	sessionIDHeader = "X-Transmission-Session-Id"
	// Below Prometheus' default scrape timeout, so a hung Transmission fails the scrape.
	requestTimeout = 10 * time.Second
)

type (
	// User to authenticate with Transmission
	User struct {
		Username string
		Password string
	}
	// Client connects to transmission via HTTP. Safe for concurrent use: the
	// collectors share one.
	Client struct {
		URL  string
		User *User

		client    http.Client
		mu        sync.Mutex
		sessionID string // guarded by mu
	}

	request struct {
		Method    string            `json:"method"`
		Arguments *requestArguments `json:"arguments,omitempty"`
	}
	requestArguments struct {
		Fields []string `json:"fields"`
	}
	response[T any] struct {
		Arguments T      `json:"arguments"`
		Result    string `json:"result"`
	}
)

// New creates a Client for the Transmission at url
func New(url string, user *User) *Client {
	return &Client{
		URL:    url + endpoint,
		User:   user,
		client: http.Client{Timeout: requestTimeout},
	}
}

// call sends req and decodes the arguments of its response.
func call[T any](c *Client, req request) (*T, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resBody, err := c.post(body)
	if err != nil {
		return nil, err
	}

	var res response[T]
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, err
	}
	if res.Result != "success" {
		return nil, fmt.Errorf("%s: %s", req.Method, res.Result)
	}
	return &res.Arguments, nil
}

func (c *Client) post(body []byte) ([]byte, error) {
	res, err := c.do(body)
	if err != nil {
		return nil, err
	}

	// 409: the session id is missing or expired, the response carries the current one.
	if res.StatusCode == http.StatusConflict {
		_ = res.Body.Close()
		c.mu.Lock()
		c.sessionID = res.Header.Get(sessionIDHeader)
		c.mu.Unlock()
		if res, err = c.do(body); err != nil {
			return nil, err
		}
	}
	defer func() { _ = res.Body.Close() }()

	switch res.StatusCode {
	case http.StatusOK:
		return io.ReadAll(res.Body)
	case http.StatusUnauthorized:
		return nil, errors.New("authorization failed, check your username and password and make sure the ip is whitelisted")
	default:
		return nil, fmt.Errorf("unexpected status %s", res.Status)
	}
}

func (c *Client) do(body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	req.Header.Set(sessionIDHeader, c.sessionID)
	c.mu.Unlock()

	if c.User != nil {
		req.SetBasicAuth(c.User.Username, c.User.Password)
	}

	return c.client.Do(req)
}

// GetTorrents get a list of torrents
func (c *Client) GetTorrents() ([]Torrent, error) {
	args, err := call[struct {
		Torrents []Torrent `json:"torrents"`
	}](c, request{Method: "torrent-get", Arguments: &requestArguments{Fields: torrentFields}})
	if err != nil {
		return nil, err
	}
	return args.Torrents, nil
}

// GetSession gets the current session from transmission
func (c *Client) GetSession() (*Session, error) {
	return call[Session](c, request{Method: "session-get"})
}

// GetSessionStats gets stats on the current & cumulative session
func (c *Client) GetSessionStats() (*SessionStats, error) {
	return call[SessionStats](c, request{Method: "session-stats"})
}
