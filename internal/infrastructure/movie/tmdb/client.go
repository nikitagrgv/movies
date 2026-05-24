package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    *url.URL
	imageURL   *url.URL
	token      string
	httpClient *http.Client
}

func NewClient(baseURL string, imageURL string, token string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	image, err := url.Parse(imageURL)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 5 * time.Second}
	return &Client{
		baseURL:    base,
		imageURL:   image,
		token:      token,
		httpClient: c,
	}, nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	full := c.baseURL.JoinPath(path)
	full.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tmdb error: status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode tmdb response: %w", err)
	}
	return nil
}

func (c *Client) getImageURL(path string) string {
	if path == "" {
		return ""
	}

	// TODO: hardcoded size 200
	path = strings.TrimPrefix(path, "/")
	full := c.imageURL.JoinPath("w200", path)
	return full.String()
}
