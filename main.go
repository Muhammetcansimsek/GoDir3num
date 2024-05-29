package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Options defines the options for the scanner
type Options struct {
	URL      string
	Wordlist []string
	Threads  int
	Verbose  bool
}

func main() {
	url := flag.String("url", "", "The base URL to fuzz")
	wordlistFile := flag.String("wordlist", "", "Path to the wordlist file")
	threads := flag.Int("threads", 150, "Number of concurrent threads (default 150)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *url == "" || *wordlistFile == "" {
		fmt.Println("URL and wordlist are required")
		flag.Usage()
		return
	}

	*url = strings.TrimRight(*url, "/")

	wordlist, err := readWordlist(*wordlistFile)
	if err != nil {
		fmt.Printf("Error reading wordlist: %v\n", err)
		return
	}

	options := Options{
		URL:      *url,
		Wordlist: wordlist,
		Threads:  *threads,
		Verbose:  *verbose,
	}

	client := createHTTPClient(*threads)
	processURLs(client, options)
}

// createHTTPClient creates a custom HTTP client
func createHTTPClient(maxConnsPerHost int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: maxConnsPerHost,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent following redirects
			return http.ErrUseLastResponse
		},
	}
}

// processURLs handles the concurrent processing of URLs
func processURLs(client *http.Client, options Options) {
	var wg sync.WaitGroup
	urls := make(chan string, options.Threads)

	for i := 0; i < options.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urls {
				checkAndPrintURL(client, url, options.Verbose, options.URL)
			}
		}()
	}

	for _, path := range options.Wordlist {
		urls <- fmt.Sprintf("%s/%s", options.URL, path)
	}
	close(urls)
	wg.Wait()
}

// readWordlist reads the wordlist from a file
func readWordlist(filePath string) ([]string, error) {
	fmt.Printf("Reading wordlist from %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open wordlist file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wordlist file: %w", err)
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("wordlist is empty")
	}
	return lines, nil
}

// checkAndPrintURL checks a single URL using the provided HTTP client and prints results
func checkAndPrintURL(client *http.Client, url string, verbose bool, baseURL string) {
	resp, err := client.Get(url)
	if err != nil {
		if verbose {
			fmt.Printf("Error fetching %s: %v\n", url, err)
		}
		return
	}
	defer resp.Body.Close()

	// Filter based on verbose flag
	if !verbose && (resp.StatusCode >= 400 && resp.StatusCode < 500) {
		return // Ignore 400-499 status codes if not in verbose mode
	}

	result := Result{
		URL:        baseURL,
		Path:       strings.TrimPrefix(url, baseURL+"/"),
		Verbose:    verbose,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
		Size:       resp.ContentLength,
		Found:      resp.StatusCode >= 200 && resp.StatusCode < 300,
	}

	output, err := result.ResultToString()
	if err != nil {
		fmt.Printf("Error converting result to string: %v\n", err)
		return
	}
	fmt.Print(output)
}
