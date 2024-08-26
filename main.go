package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func generateWordSet() map[string]struct{} {
	wordSet := map[string]struct{}{}

	// Read the file into a byte slice
	fileContent, err := os.ReadFile("assets/word-bank.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Scan each word and add it to the set
	reader := bytes.NewReader(fileContent)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		word := scanner.Text()
		// We use an empty struct as the value to mimic a set
		wordSet[word] = struct{}{}
	}
	return wordSet
}

func main() {

	wordSet := generateWordSet()
	// Example: Check if a specific word is in the set
	checkWord := "exampleeee"
	// _ is the value, exists is the boolean
	if _, exists := wordSet[checkWord]; exists {
		fmt.Printf("The word '%s' exists in the set.\n", checkWord)
	} else {
		fmt.Printf("The word '%s' does not exist in the set.\n", checkWord)
	}
}
