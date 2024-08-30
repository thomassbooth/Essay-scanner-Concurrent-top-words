package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// this removes all punctuation from the essay so we can split
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9' ]+`)

func cleanEssay(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func fetchAndProcessEssay(url string) (string, error, bool) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s: %v", url, err), false
	}
	defer response.Body.Close()

	// if our status code is 999 this means weve been rate limited, so retry
	if response.StatusCode == 999 {
		return "", fmt.Errorf("received 999 error for %s", url), true
		// other status codes that arent 200 are errors and we dont wanna retry
	} else if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response %d for %s", response.StatusCode, url), false
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of %s: %v", url, err), false
	}

	// we have got our raw essay text, now we need to clean it up
	text := ParseHTMLFile(string(body))
	cleanedEssay := cleanEssay(text)
	// fmt.Println(cleanedEssay)
	return cleanedEssay, nil, false
}

// Extract text from HTML, focusing on <div class="caas-body">
func ParseHTMLFile(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}
	return TargetDivDFS(doc)
}

// Traverse the HTML nodes to find the <div class="caas-body"> and extract text from its <p> tags
func TargetDivDFS(node *html.Node) string {
	var result string

	// Check if this node is the target div
	if node.Data == "div" {
		// attributes is a slice of key-value pairs, we need to loop through them untill we find the class and check its value
		for _, attr := range node.Attr {
			// check our div has the class were looking for
			if attr.Key == "class" && attr.Val == "caas-body" {
				// Extract text from <p> tags within this div, ignoring other nested divs
				result += FindPTags(node)
				break
			}
		}
	}

	// Continue to traverse the document tree
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result += TargetDivDFS(c)
	}
	return result
}

// Extract text from <p> tags within the node, ignoring non-<p> elements
func FindPTags(node *html.Node) string {
	var result string

	// weve passed in the div node, we need to loop through its childen to find the <p> tags since they are directly below the div
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "div" {
			// Ignore any nested divs
			continue
		}
		// If it's a <p> tag, extract the text
		if c.Data == "p" {
			result += PTagDFS(c) // Add a newline to separate paragraphs
		}
	}
	return result
}

// Traverse the HTML nodes to extract all text
func PTagDFS(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var result string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result += PTagDFS(c)
	}

	return result
}
