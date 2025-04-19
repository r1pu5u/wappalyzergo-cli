package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

func main() {
	// Command-line flags
	listFile := flag.String("l", "", "Path to file containing list of URLs")
	threads := flag.Int("t", 5, "Number of concurrent threads (default 5)")
	output := flag.String("o", "", "Output JSON file path")
	flag.Parse()

	if *listFile == "" || *output == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Read URLs from the provided list file
	file, err := os.Open(*listFile)
	if err != nil {
		log.Fatalf("Failed to open list file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var urls []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading list file: %v", err)
	}

	// Prepare concurrency primitives
	sem := make(chan struct{}, *threads)
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string]map[string]struct{})

	wa, err := wappalyzer.New()
	if err != nil {
		log.Fatalf("Failed to initialize Wappalyzer: %v", err)
	}

	// Launch scanning goroutines
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			sem <- struct{}{}            // Acquire slot
			defer func() { <-sem }()      // Release slot

			fmt.Printf("Processing %s\n", u)

			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				log.Printf("Failed to create request for %s: %v", u, err)
				return
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

			client := http.DefaultClient
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Failed to GET %s: %v", u, err)
				return
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response from %s: %v", u, err)
				return
			}

			// Fingerprint technologies
			fps := wa.Fingerprint(resp.Header, data)

			// Store results
			mu.Lock()
			results[u] = fps
			mu.Unlock()
		}(url)
	}

	// Wait for all scans to complete
	wg.Wait()

	// Marshal to JSON
	outData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal results: %v", err)
	}

	// Write to output file
	err = os.WriteFile(*output, outData, 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Scan complete. Results written to %s\n", *output)
}
