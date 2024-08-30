package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type BackOff struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxElapsedTime  time.Duration
	FirstCallDelay  time.Duration
}

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

// applies our settings we set in main for our exponential backoff
func SetupBackoff(b BackOff) *backoff.ExponentialBackOff {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = b.InitialInterval
	expBackoff.MaxInterval = b.MaxInterval
	expBackoff.Multiplier = b.Multiplier
	expBackoff.MaxElapsedTime = b.MaxElapsedTime

	return expBackoff
}

// https://stackoverflow.com/questions/27117896/how-to-pretty-print-variables
// PrettyPrint prints a slice of WordCount pairs in a pretty JSON format.
func PrettyPrint(pairs []WordCount) {
	// Convert the slice of key-value pairs to a map
	pairMap := make(map[string]string)
	for _, pair := range pairs {
		pairMap[pair.Word] = strconv.Itoa(pair.Count)
	}

	// Marshal the map into a pretty JSON format
	jsonBytes, err := json.MarshalIndent(pairMap, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the pretty JSON string
	fmt.Println(string(jsonBytes))
}
