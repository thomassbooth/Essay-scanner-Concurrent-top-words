# Essay Scanner - Top K Words

## Overview
- a Go application designed to process essays fetched from https://www.engadget.com. It parses HTML content, counts word frequencies, and maintains the top k most frequent words using a heap. The application supports concurrent processing and implements rate limiting with exponential backoff.

## Project Structure

- **`main.go`**
  - Entry point of the application, runs the application and holds constants for adjusting efficiency within the application.

- **`essays.go`**
  - Parses the HTML returned from fetch requests.
  - Fetches all essays and processes them.
  - Cleans parsed data.

- **`minheap.go`**
  - Implements a min-heap structure for storing key-value pairs.

- **`wordcounter.go`**
  - Contains the `topKelements` algorithm, which maintains a heap of size `k` where it keeps the highest frequency elements and replaces the smallest one as needed.

- **`workerpool.go`**
  - Manages concurrency with a pool of workers.
  - Each worker uses exponential backoff to handle retries.

- **`utils.go`**
  - Provides generic utility functions, such as file parsing and pretty printing.

## Running the application

1. **Install dependencies**

    Install dependencies needed for the project

    ```bash
    go mod tidy
    ```

2. **Run the application**

    Ensure youre at the root directory of the project

    ```bash
    go run .
    ```

## Additions

- **Constants and Variables**:
    (these are found in main.go)
    - **`wordSetPath`**: Path to the file containing the word bank. Default is `"assets/word-bank.txt"`.
    - **`essayUrlsPath`**: Path to the file containing the list of essay URLs. Default is `"assets/test-urls.txt"`.
    - **`workers`**: Number of concurrent workers used for processing essays. Default is `4`.
    - **`maxWords`**: Maximum number of top frequent words to retain. Default is `10`.
    - **`minWordLen`**: Minimum length of words to consider for counting. Default is `3`.
    - **`backOffIntervals`**: Configuration for exponential backoff used for handling retries.
      - **`InitialInterval`**: Initial wait time before retrying. Default is `10 * time.Second`.
      - **`MaxInterval`**: Maximum wait time before retrying. Default is `5 * time.Second`.
      - **`Multiplier`**: Factor by which the wait time increases. Default is `2.0`.
      - **`MaxElapsedTime`**: Maximum total time for which retries will be attempted. Default is `30 * time.Second`.
      - **`FirstCallDelay`**: Delay before the first retry attempt. Default is `500 * time.Millisecond`.

