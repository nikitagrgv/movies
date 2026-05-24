package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		ForceAttemptHTTP2:   true,
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}
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

	var lastErr error
	for i := 0; i < 3; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, full.String(), nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
			continue
		}

		func() {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				lastErr = fmt.Errorf("tmdb error: status=%d body=%s", resp.StatusCode, string(body))
				return
			}

			lastErr = json.NewDecoder(resp.Body).Decode(out)
		}()

		if lastErr == nil {
			return nil
		}

		if resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
			continue
		}

		return lastErr
	}

	return lastErr
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
