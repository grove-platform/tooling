package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// Define command-line flags
	csvFile := flag.String("csv", "", "Path to CSV file containing URLs")
	outputDir := flag.String("output", "markdown-output", "Output directory for markdown files")
	workers := flag.Int("workers", 10, "Number of concurrent download workers")
	rateLimit := flag.Float64("rate-limit", 5.0, "Maximum requests per second (0 for unlimited)")
	flag.Parse()

	if *csvFile == "" {
		fmt.Println("Error: -csv flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Create rate limiter
	var limiter *rate.Limiter
	if *rateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(*rateLimit), 1)
		fmt.Printf("Rate limiting enabled: %.1f requests/second\n", *rateLimit)
	} else {
		limiter = rate.NewLimiter(rate.Inf, 0)
		fmt.Println("Rate limiting disabled")
	}

	// Read and process CSV file
	if err := processCSV(*csvFile, *outputDir, *workers, limiter); err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDownload complete!")
}

// downloadJob represents a download task
type downloadJob struct {
	url     string
	lineNum int
}

// downloadResult represents the result of a download
type downloadResult struct {
	url     string
	lineNum int
	err     error
}

// isHeaderRow detects if a CSV row is a header row
func isHeaderRow(record []string) bool {
	if len(record) < 2 {
		return false
	}

	// Check for common header patterns from create-url-list output
	// Output headers: "Rank", "Page", "Number of Page Views"
	// Input headers: "Page", "Page Subsite", "Measure Names", "Measure Values", "Min. Aux"

	// Check first and second columns
	firstCol := strings.ToLower(strings.TrimSpace(record[0]))
	secondCol := strings.ToLower(strings.TrimSpace(record[1]))

	// If either column looks like a URL (contains mongodb.com or www.), it's not a header
	if strings.Contains(firstCol, "mongodb.com") || strings.Contains(firstCol, "www.") {
		return false
	}
	if strings.Contains(secondCol, "mongodb.com") || strings.Contains(secondCol, "www.") {
		return false
	}

	// Common header patterns in first column
	firstColKeywords := []string{"rank", "page"}
	for _, keyword := range firstColKeywords {
		if firstCol == keyword {
			return true
		}
	}

	// Common header patterns in second column
	secondColKeywords := []string{"page", "measure", "url"}
	for _, keyword := range secondColKeywords {
		if strings.Contains(secondCol, keyword) {
			return true
		}
	}

	// If the second column doesn't contain a dot, it's likely not a URL
	if !strings.Contains(secondCol, ".") {
		return true
	}

	return false
}

// processCSV reads the CSV file and downloads markdown for each URL using concurrent workers
func processCSV(csvPath, outputDir string, numWorkers int, limiter *rate.Limiter) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Read all URLs from CSV first
	reader := csv.NewReader(file)
	var jobs []downloadJob
	lineNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV at line %d: %w", lineNum, err)
		}

		lineNum++

		// Skip header row if present
		// Common headers: "Rank,Page" or "Rank,Page,Number of Page Views" or "Page,Page Subsite,..."
		if lineNum == 1 && isHeaderRow(record) {
			fmt.Println("Detected and skipping header row")
			continue
		}

		// Expect format: rank,url,metric (or just url,metric)
		if len(record) < 2 {
			fmt.Printf("Warning: Line %d has insufficient columns, skipping\n", lineNum)
			continue
		}

		jobs = append(jobs, downloadJob{
			url:     record[1],
			lineNum: lineNum,
		})
	}

	if len(jobs) == 0 {
		fmt.Println("No URLs to process")
		return nil
	}

	fmt.Printf("Processing %d URLs with %d workers...\n", len(jobs), numWorkers)

	// Create channels for jobs and results
	jobChan := make(chan downloadJob, len(jobs))
	resultChan := make(chan downloadResult, len(jobs))

	// Start worker goroutines
	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i+1, jobChan, resultChan, outputDir, limiter, &wg)
	}

	// Send jobs to workers
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	successCount := 0
	errorCount := 0
	for result := range resultChan {
		if result.err != nil {
			fmt.Printf("Error downloading %s: %v\n", result.url, result.err)
			errorCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("\nProcessed %d URLs: %d successful, %d errors\n", len(jobs), successCount, errorCount)
	return nil
}

// worker processes download jobs from the job channel
func worker(ctx context.Context, id int, jobs <-chan downloadJob, results chan<- downloadResult, outputDir string, limiter *rate.Limiter, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Wait for rate limiter before making request
		if err := limiter.Wait(ctx); err != nil {
			results <- downloadResult{
				url:     job.url,
				lineNum: job.lineNum,
				err:     fmt.Errorf("rate limiter error: %w", err),
			}
			continue
		}

		err := downloadMarkdown(job.url, outputDir)
		results <- downloadResult{
			url:     job.url,
			lineNum: job.lineNum,
			err:     err,
		}
	}
}

// downloadMarkdown downloads the markdown version of a documentation page
func downloadMarkdown(pageURL, outputDir string) error {
	// Extract the page slug and construct markdown URL
	slug, err := extractPageSlug(pageURL)
	if err != nil {
		return fmt.Errorf("failed to extract page slug: %w", err)
	}

	// Check if file already exists
	outputPath := filepath.Join(outputDir, slug+".md")
	if _, err := os.Stat(outputPath); err == nil {
		// File already exists, skip download
		return nil
	}

	mdURL, err := constructMarkdownURL(pageURL)
	if err != nil {
		return fmt.Errorf("failed to construct markdown URL: %w", err)
	}

	// Create HTTP request with User-Agent header
	req, err := http.NewRequest("GET", mdURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; get-docs-markdown/1.0)")

	// Download the markdown content
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Create subdirectories if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create subdirectories: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Write content to file
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// extractPageSlug extracts the page slug from a URL (everything after www.mongodb.com/docs/)
// Includes language and version prefixes to ensure uniqueness
func extractPageSlug(pageURL string) (string, error) {
	// Remove protocol if present
	pageURL = strings.TrimPrefix(pageURL, "http://")
	pageURL = strings.TrimPrefix(pageURL, "https://")

	// Find the position of www.mongodb.com/
	domainIndex := strings.Index(pageURL, "www.mongodb.com/")
	if domainIndex == -1 {
		return "", fmt.Errorf("URL does not contain 'www.mongodb.com/': %s", pageURL)
	}

	// Extract everything after www.mongodb.com/
	afterDomain := pageURL[domainIndex+16:] // +16 to skip "www.mongodb.com/"

	// Check if there's a language/version prefix before /docs/
	docsIndex := strings.Index(afterDomain, "docs/")
	if docsIndex == -1 {
		return "", fmt.Errorf("URL does not contain 'docs/': %s", pageURL)
	}

	// Include any language/version prefix before docs/
	var slug string
	if docsIndex > 0 {
		// There's a prefix (e.g., "zh-cn/" or "v7.0/")
		prefix := afterDomain[:docsIndex]
		prefix = strings.Trim(prefix, "/")
		afterDocs := afterDomain[docsIndex+5:] // +5 to skip "docs/"
		slug = prefix + "/" + afterDocs
	} else {
		// No prefix, just extract after docs/
		slug = afterDomain[docsIndex+5:] // +5 to skip "docs/"
	}

	// Remove query parameters and anchors first
	if idx := strings.IndexAny(slug, "?#"); idx != -1 {
		slug = slug[:idx]
	}

	// Then remove trailing slash
	slug = strings.TrimSuffix(slug, "/")

	// If slug is empty, use "index"
	if slug == "" {
		slug = "index"
	}

	return slug, nil
}

// constructMarkdownURL creates the markdown URL from a page URL
func constructMarkdownURL(pageURL string) (string, error) {
	// Add https:// if not present
	if !strings.HasPrefix(pageURL, "http://") && !strings.HasPrefix(pageURL, "https://") {
		pageURL = "https://" + pageURL
	}

	// Parse the URL
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Remove query parameters and fragment
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""

	// Remove trailing slash from path
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	// Add .md extension
	parsedURL.Path += ".md"

	return parsedURL.String(), nil
}
