package main

import (
	"fmt"
)

const wordSetPath = "assets/word-bank.txt"
const essayUrlsPath = "assets/endg-urls.txt"

func main() {
	wordSet := GenerateWordSet(wordSetPath)
	urls := FormatEssayUrls(essayUrlsPath)
	// urls := []string{"https://www.engadget.com/2019/08/24/trump-tries-to-overturn-ruling-stopping-him-from-blocking-twitte/"}
	processEssays(urls)
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
