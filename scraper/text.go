package scraper

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Extracts all visible text from an HTML page
func ExtractTextFromURL(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch page")
	}

	// Fallback to manual DOM traversal
	return extractTextFromHTML(resp.Body)
}

func extractTextFromHTML(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode && !isIgnorable(n.Parent) {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				b.WriteString(text + " ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	return strings.TrimSpace(b.String()), nil
}

func isIgnorable(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode {
		return false
	}

	switch n.Data {
	case "script", "style", "head", "noscript":
		return true
	}
	return false
}
