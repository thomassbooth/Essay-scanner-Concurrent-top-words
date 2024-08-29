package main

import (
	"strings"

	"golang.org/x/net/html"
)

// Function to extract text from HTML, focusing on <div class="caas-body">
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

// Count words in the text
func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}
