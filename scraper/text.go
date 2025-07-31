package scraper

import (
	"errors"
	"net/http"
	"net/url"

	"time"

	"github.com/go-shiori/go-readability"
)

func ExtractTextFromURL(u string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch the page")
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", errors.New("failed to parse URL")
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", err
	}

	if len(article.TextContent) < 100 {
		return "", errors.New("extract failed: content too short")
	}

	return article.TextContent, nil
}
