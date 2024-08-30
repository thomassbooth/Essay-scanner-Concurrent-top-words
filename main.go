package main

const wordSetPath = "assets/word-bank.txt"
const essayUrlsPath = "assets/test-urls.txt"

func main() {
	urls := FormatEssayUrls(essayUrlsPath)
	processEssays(urls)
}
