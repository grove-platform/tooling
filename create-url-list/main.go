package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Record struct {
	Page          string
	MeasureValues int
}

type Config struct {
	IgnoreURLs    []string `yaml:"ignore_urls"`
	ShowPageviews bool     `yaml:"show_pageviews"`
	ShowHeaders   bool     `yaml:"show_headers"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s [--quiet] [--contains <substring>] <csv-file-path> [range] [output-path]", os.Args[0])
	}

	// Check for --quiet and --contains flags
	quiet := false
	containsFilter := ""
	args := os.Args[1:]

	// Process flags
	for len(args) > 0 && strings.HasPrefix(args[0], "--") {
		if args[0] == "--quiet" {
			quiet = true
			args = args[1:] // Remove --quiet from args
		} else if args[0] == "--contains" {
			if len(args) < 2 {
				return fmt.Errorf("--contains flag requires a substring argument")
			}
			containsFilter = args[1]
			args = args[2:] // Remove --contains and its argument from args
		} else {
			return fmt.Errorf("unknown flag: %s", args[0])
		}
	}

	if len(args) < 1 {
		return fmt.Errorf("usage: %s [--quiet] [--contains <substring>] <csv-file-path> [range] [output-path]", os.Args[0])
	}

	inputPath := args[0]
	rangeStr := "1-250" // default range
	outputPath := ""

	if len(args) >= 2 {
		rangeStr = args[1]
	}
	if len(args) >= 3 {
		outputPath = args[2]
	}

	// Parse range
	minVal, maxVal, err := parseRange(rangeStr)
	if err != nil {
		return fmt.Errorf("invalid range format: %v", err)
	}

	// Validate input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Load config (optional)
	config, err := loadConfig("config.yml")
	if err != nil {
		// Config is optional, so we just log a warning if it fails
		if !quiet {
			fmt.Fprintf(os.Stderr, "Warning: Could not load config.yml: %v\n", err)
		}
		config = &Config{} // Use empty config
	}

	// Read and process CSV
	records, err := processCSV(inputPath, config.IgnoreURLs, containsFilter, quiet)
	if err != nil {
		return err
	}

	// Generate output
	outputFilePath, err := writeOutput(records, outputPath, rangeStr, minVal, maxVal, config.ShowPageviews, config.ShowHeaders)
	if err != nil {
		return err
	}

	// Print success message
	if !quiet {
		fmt.Printf("Successfully parsed input file `%s` and created output file at `%s`\n", inputPath, outputFilePath)
	}

	return nil
}

func parseRange(rangeStr string) (int, int, error) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("range must be in format 'min-max'")
	}

	min, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minimum value: %v", err)
	}

	max, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid maximum value: %v", err)
	}

	if min > max {
		return 0, 0, fmt.Errorf("minimum value cannot be greater than maximum value")
	}

	return min, max, nil
}

func loadConfig(configPath string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func processCSV(inputPath string, ignoreURLs []string, containsFilter string, quiet bool) ([]Record, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	// Find column indices
	pageIdx, measureNamesIdx, measureValuesIdx := -1, -1, -1
	for i, col := range header {
		switch col {
		case "Page":
			pageIdx = i
		case "Measure Names":
			measureNamesIdx = i
		case "Measure Values":
			measureValuesIdx = i
		}
	}

	// Validate required columns exist
	if pageIdx == -1 || measureNamesIdx == -1 || measureValuesIdx == -1 {
		return nil, fmt.Errorf("missing required columns (Page, Measure Names, Measure Values)")
	}

	// Create a map for fast lookup of ignored URLs
	ignoreMap := make(map[string]bool)
	for _, url := range ignoreURLs {
		ignoreMap[url] = true
	}

	// Read and collect all Pageviews records
	var records []Record
	var skippedURLs []string
	var ignoredURLs []string
	var filteredURLs []string
	for {
		row, err := reader.Read()
		if err != nil {
			break // EOF or error
		}

		// Skip if not enough columns
		if len(row) <= pageIdx || len(row) <= measureNamesIdx || len(row) <= measureValuesIdx {
			continue
		}

		// Filter by Measure Names = "Pageviews"
		if row[measureNamesIdx] != "Pageviews" {
			continue
		}

		// Validate URL structure
		page := row[pageIdx]
		if !strings.HasPrefix(page, "www.") {
			skippedURLs = append(skippedURLs, page)
			continue
		}

		// Check if URL should be ignored
		if ignoreMap[page] {
			ignoredURLs = append(ignoredURLs, page)
			continue
		}

		// Filter by contains substring if specified
		if containsFilter != "" && !strings.Contains(page, containsFilter) {
			filteredURLs = append(filteredURLs, page)
			continue
		}

		// Parse Measure Values
		measureValue, err := strconv.Atoi(row[measureValuesIdx])
		if err != nil {
			continue // Skip non-integer values
		}

		records = append(records, Record{
			Page:          page,
			MeasureValues: measureValue,
		})
	}

	// Report skipped URLs
	if !quiet && len(skippedURLs) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: Skipped %d URL(s) that do not match expected structure (www.*):\n", len(skippedURLs))
		for _, url := range skippedURLs {
			fmt.Fprintf(os.Stderr, "  - %s\n", url)
		}
	}

	// Report ignored URLs
	if !quiet && len(ignoredURLs) > 0 {
		fmt.Fprintf(os.Stderr, "Info: Ignored %d URL(s) from config:\n", len(ignoredURLs))
		for _, url := range ignoredURLs {
			fmt.Fprintf(os.Stderr, "  - %s\n", url)
		}
	}

	// Report filtered URLs
	if !quiet && len(filteredURLs) > 0 {
		fmt.Fprintf(os.Stderr, "Info: Filtered out %d URL(s) not containing '%s':\n", len(filteredURLs), containsFilter)
		for _, url := range filteredURLs {
			fmt.Fprintf(os.Stderr, "  - %s\n", url)
		}
	}

	return records, nil
}

func writeOutput(records []Record, outputPath, rangeStr string, minRank, maxRank int, showPageviews, showHeaders bool) (string, error) {
	// Sort by Measure Values (highest to lowest) to establish ranking
	sort.Slice(records, func(i, j int) bool {
		return records[i].MeasureValues > records[j].MeasureValues
	})

	// Slice to get only the entries within the specified rank range
	// minRank and maxRank are 1-based, so we need to convert to 0-based indices
	startIdx := minRank - 1
	endIdx := maxRank

	// Ensure we don't go out of bounds
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(records) {
		endIdx = len(records)
	}
	if startIdx >= len(records) {
		// No records in this range
		records = []Record{}
	} else {
		records = records[startIdx:endIdx]
	}

	// Determine output directory and filename
	var outputDir, filename string
	if outputPath != "" {
		outputDir = filepath.Dir(outputPath)
		filename = filepath.Base(outputPath)
	} else {
		outputDir = "output"
		// Generate filename: YYYY-MM-DD_HH-MM-SS_range.csv
		now := time.Now()
		filename = fmt.Sprintf("%s_%s.csv",
			now.Format("2006-01-02_15-04-05"),
			rangeStr)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create output file
	outputFilePath := filepath.Join(outputDir, filename)
	file, err := os.Create(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers if enabled
	if showHeaders {
		var headers []string
		if showPageviews {
			headers = []string{"Rank", "Page", "Number of Page Views"}
		} else {
			headers = []string{"Rank", "Page"}
		}
		if err := writer.Write(headers); err != nil {
			return "", fmt.Errorf("failed to write headers: %v", err)
		}
	}

	// Write records with rank number, URL, and optionally pageviews
	for i, record := range records {
		rank := startIdx + i + 1 // Calculate the actual rank
		var row []string
		if showPageviews {
			row = []string{
				strconv.Itoa(rank),
				record.Page,
				strconv.Itoa(record.MeasureValues),
			}
		} else {
			row = []string{
				strconv.Itoa(rank),
				record.Page,
			}
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write record: %v", err)
		}
	}

	return outputFilePath, nil
}
