package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

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
