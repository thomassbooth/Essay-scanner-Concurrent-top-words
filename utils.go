package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

// takes in a file path, reads the file and processes each line by the function passed in
func processFile(filePath string, processLine func(string)) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	// Scan each line and process it
	reader := bytes.NewReader(fileContent)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		processLine(scanner.Text())
	}
}

// takes in a file path and returns a slice of urls
func FormatEssayUrls(path string) []string {
	var urls []string
	processFile(path, func(line string) {
		urls = append(urls, line)
	})
	return urls
}

// takes in a file path and returns a set of words
func GenerateWordSet(path string) map[string]struct{} {
	wordSet := map[string]struct{}{}
	// Read the file into a byte slice
	processFile(path, func(word string) {
		wordSet[word] = struct{}{}
	})
	return wordSet
}
