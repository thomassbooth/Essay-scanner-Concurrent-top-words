package main

import "time"

const wordSetPath = "assets/word-bank.txt"
const essayUrlsPath = "assets/test-urls.txt"

const workers = 4
const maxWords = 10
const minWordLen = 3

var backOffIntervals = BackOff{
	InitialInterval: 500 * time.Millisecond,
	MaxInterval:     5 * time.Second,
	Multiplier:      2.0,
	MaxElapsedTime:  30 * time.Second,
}

func main() {
	urls := FormatEssayUrls(essayUrlsPath)
	results := processEssays(urls)
	PrettyPrint(results)
}
