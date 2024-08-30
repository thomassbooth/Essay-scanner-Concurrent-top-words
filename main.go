package main

import (
	"fmt"
	"time"
)

const wordSetPath = "assets/word-bank.txt"
const essayUrlsPath = "assets/test-urls.txt"

const workers = 4
const maxWords = 10
const minWordLen = 3

var backOffIntervals = BackOff{
	InitialInterval: 10 * time.Second,
	MaxInterval:     5 * time.Second,
	Multiplier:      2.0,
	MaxElapsedTime:  30 * time.Second,
	FirstCallDelay:  500 * time.Millisecond,
}

func main() {
	urls := FormatEssayUrls(essayUrlsPath)
	results := processEssays(urls)
	fmt.Println("Processing Completed, Top 10 words:")
	PrettyPrint(results)
}
