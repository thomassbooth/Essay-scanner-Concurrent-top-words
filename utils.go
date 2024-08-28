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
