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
)

type Record struct {
	Page          string
	MeasureValues int
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
		return fmt.Errorf("usage: %s <csv-file-path> [range] [output-path]", os.Args[0])
	}

	inputPath := os.Args[1]
	rangeStr := "1-250" // default range
	outputPath := ""

	if len(os.Args) >= 3 {
		rangeStr = os.Args[2]
	}
	if len(os.Args) >= 4 {
		outputPath = os.Args[3]
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

	// Read and process CSV
	records, err := processCSV(inputPath)
	if err != nil {
		return err
	}

	// Generate output
	if err := writeOutput(records, outputPath, rangeStr, minVal, maxVal); err != nil {
		return err
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

func processCSV(inputPath string) ([]Record, error) {
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

	// Read and collect all Pageviews records
	var records []Record
	var skippedURLs []string
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
	if len(skippedURLs) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: Skipped %d URL(s) that do not match expected structure (www.*):\n", len(skippedURLs))
		for _, url := range skippedURLs {
			fmt.Fprintf(os.Stderr, "  - %s\n", url)
		}
	}

	return records, nil
}

func writeOutput(records []Record, outputPath, rangeStr string, minRank, maxRank int) error {
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
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create output file
	outputFilePath := filepath.Join(outputDir, filename)
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write records with rank number and URL
	for i, record := range records {
		rank := startIdx + i + 1 // Calculate the actual rank
		if err := writer.Write([]string{
			strconv.Itoa(rank),
			record.Page,
		}); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}
