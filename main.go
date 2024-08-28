package main

import (
	"fmt"
)

func FormatEssayUrls(path string) []string {
	var urls []string
	processFile(path, func(line string) {
		urls = append(urls, line)
	})
	return urls
}

func GenerateWordSet(path string) map[string]struct{} {
	wordSet := map[string]struct{}{}
	// Read the file into a byte slice
	processFile(path, func(word string) {
		wordSet[word] = struct{}{}
	})
	return wordSet
}

func main() {
	wordSet := GenerateWordSet("assets/word-bank.txt")
	// urls := FormatEssayUrls("assets/endg-urls.txt")
	// processEssays(urls)
	// for _, url := range urls {
	// 	fmt.Println(url)
	// }
	// Example: Check if a specific word is in the set
	checkWord := "exampleeee"
	// _ is the value, exists is the boolean
	if _, exists := wordSet[checkWord]; exists {
		fmt.Printf("The word '%s' exists in the set.\n", checkWord)
	} else {
		fmt.Printf("The word '%s' does not exist in the set.\n", checkWord)
	}
}
