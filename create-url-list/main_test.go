package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseRange tests the parseRange function with various inputs
func TestParseRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMin   int
		wantMax   int
		wantError bool
	}{
		{"valid range", "1-250", 1, 250, false},
		{"valid range with spaces", "10 - 100", 10, 100, false},
		{"single digit range", "5-9", 5, 9, false},
		{"large range", "1-1000000", 1, 1000000, false},
		{"invalid format - no dash", "100", 0, 0, true},
		{"invalid format - multiple dashes", "1-2-3", 0, 0, true},
		{"invalid min - not a number", "abc-100", 0, 0, true},
		{"invalid max - not a number", "1-xyz", 0, 0, true},
		{"min greater than max", "100-50", 0, 0, true},
		{"empty string", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max, err := parseRange(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("parseRange(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseRange(%q) unexpected error: %v", tt.input, err)
				}
				if min != tt.wantMin {
					t.Errorf("parseRange(%q) min = %d, want %d", tt.input, min, tt.wantMin)
				}
				if max != tt.wantMax {
					t.Errorf("parseRange(%q) max = %d, want %d", tt.input, max, tt.wantMax)
				}
			}
		})
	}
}

// TestProcessCSV_MissingColumns tests that processCSV returns an error when required columns are missing
func TestProcessCSV_MissingColumns(t *testing.T) {
	_, err := processCSV("testdata/missing-columns.csv")
	if err == nil {
		t.Error("processCSV() expected error for missing columns, got nil")
	}
	expectedMsg := "missing required columns"
	if err != nil && !contains(err.Error(), expectedMsg) {
		t.Errorf("processCSV() error = %v, want error containing %q", err, expectedMsg)
	}
}

// TestProcessCSV_InvalidURL tests that processCSV skips URLs that don't start with www.
func TestProcessCSV_InvalidURL(t *testing.T) {
	records, err := processCSV("testdata/invalid-url.csv")
	if err != nil {
		t.Errorf("processCSV() unexpected error: %v", err)
	}
	// Should return 0 records since the only URL doesn't start with www.
	if len(records) != 0 {
		t.Errorf("processCSV() got %d records, want 0 (invalid URL should be skipped)", len(records))
	}
}

// TestProcessCSV_ValidFiltering tests that processCSV correctly collects all Pageviews records
func TestProcessCSV_ValidFiltering(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		expectedCount int
	}{
		{
			name:          "valid-with-filtering.csv collects all Pageviews",
			file:          "testdata/valid-with-filtering.csv",
			expectedCount: 6, // 50, 200, 300, 100, 1, 250 (excludes Sessions row)
		},
		{
			name:          "simple.csv with one Pageview",
			file:          "testdata/simple.csv",
			expectedCount: 1, // One Pageviews row in simple.csv
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := processCSV(tt.file)
			if err != nil {
				t.Fatalf("processCSV() unexpected error: %v", err)
			}
			if len(records) != tt.expectedCount {
				t.Errorf("processCSV() got %d records, want %d", len(records), tt.expectedCount)
			}
			// Verify all records have Pageviews (MeasureValues should be integers)
			for _, record := range records {
				if record.MeasureValues < 0 {
					t.Errorf("Record has invalid pageview value: %d", record.MeasureValues)
				}
			}
		})
	}
}

// TestProcessCSV_EmptyFile tests that processCSV handles empty CSV files
func TestProcessCSV_EmptyFile(t *testing.T) {
	records, err := processCSV("testdata/empty.csv")
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("processCSV() got %d records, want 0", len(records))
	}
}

// TestProcessCSV_FileNotFound tests that processCSV returns an error for non-existent files
func TestProcessCSV_FileNotFound(t *testing.T) {
	_, err := processCSV("testdata/nonexistent.csv")
	if err == nil {
		t.Error("processCSV() expected error for non-existent file, got nil")
	}
}

// TestWriteOutput tests the writeOutput function
func TestWriteOutput(t *testing.T) {
	tests := []struct {
		name       string
		records    []Record
		outputPath string
		rangeStr   string
		minRank    int
		maxRank    int
		wantErr    bool
	}{
		{
			name: "write top 3 records",
			records: []Record{
				{Page: "www.example.com/page1", MeasureValues: 100},
				{Page: "www.example.com/page2", MeasureValues: 50},
				{Page: "www.example.com/page3", MeasureValues: 200},
			},
			outputPath: "",
			rangeStr:   "1-3",
			minRank:    1,
			maxRank:    3,
			wantErr:    false,
		},
		{
			name: "write single record",
			records: []Record{
				{Page: "www.example.com/page1", MeasureValues: 100},
			},
			outputPath: "test-output/custom.csv",
			rangeStr:   "1-1",
			minRank:    1,
			maxRank:    1,
			wantErr:    false,
		},
		{
			name:       "write empty records",
			records:    []Record{},
			outputPath: "test-output/empty.csv",
			rangeStr:   "1-100",
			minRank:    1,
			maxRank:    100,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeOutput(tt.records, tt.outputPath, tt.rangeStr, tt.minRank, tt.maxRank)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the file was created
			if tt.outputPath == "" {
				// For default output, check that the output directory exists
				if _, err := os.Stat("output"); os.IsNotExist(err) {
					t.Error("writeOutput() did not create output directory")
				}
			}

			// Clean up test output
			if tt.outputPath != "" {
				os.RemoveAll(filepath.Dir(tt.outputPath))
			}
		})
	}

	// Clean up default output directory
	os.RemoveAll("output")
}

// TestWriteOutput_Sorting tests that records are sorted correctly by rank (highest pageviews first)
func TestWriteOutput_Sorting(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "sorted.csv")

	records := []Record{
		{Page: "www.example.com/page3", MeasureValues: 300},
		{Page: "www.example.com/page1", MeasureValues: 100},
		{Page: "www.example.com/page2", MeasureValues: 200},
		{Page: "www.example.com/page4", MeasureValues: 50},
	}

	err := writeOutput(records, outputPath, "1-4", 1, 4)
	if err != nil {
		t.Fatalf("writeOutput() unexpected error: %v", err)
	}

	// Read the output file and verify sorting
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := splitLines(string(content))
	// Should have 4 lines (no header)
	if len(lines) < 4 {
		t.Fatalf("Expected at least 4 lines, got %d", len(lines))
	}

	// Verify first line has rank 1 (highest pageviews = 300)
	if !contains(lines[0], "1,www.example.com/page3") {
		t.Errorf("First line should contain rank 1 and page3 URL, got: %s", lines[0])
	}

	// Verify last data line has rank 4 (lowest pageviews = 50)
	if !contains(lines[3], "4,www.example.com/page4") {
		t.Errorf("Last line should contain rank 4 and page4 URL, got: %s", lines[3])
	}
}

// TestWriteOutput_NoHeaders tests that output CSV has no headers
func TestWriteOutput_NoHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "no-headers.csv")

	records := []Record{
		{Page: "www.example.com/page1", MeasureValues: 100},
	}

	err := writeOutput(records, outputPath, "1-1", 1, 1)
	if err != nil {
		t.Fatalf("writeOutput() unexpected error: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := splitLines(string(content))
	// First line should be data, not headers
	if contains(lines[0], "Rank") || contains(lines[0], "Page") || contains(lines[0], "URL") {
		t.Error("Output file should not contain headers")
	}

	// First line should contain the actual data (rank and URL)
	if !contains(lines[0], "1,www.example.com/page1") {
		t.Errorf("First line should contain rank and URL, got: %s", lines[0])
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Helper function to split content into lines
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if i > start {
				lines = append(lines, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// TestProcessCSV_OnlyPageviewsFiltered tests that only Pageviews rows are included
func TestProcessCSV_OnlyPageviewsFiltered(t *testing.T) {
	records, err := processCSV("testdata/more-data.csv")
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}

	// more-data.csv has many rows with different Measure Names, but only one Pageviews row
	if len(records) != 1 {
		t.Errorf("processCSV() got %d records, want 1 (one Pageviews row)", len(records))
	}

	// Verify it's the correct Pageviews entry
	if len(records) > 0 && records[0].MeasureValues != 311105 {
		t.Errorf("Expected pageviews of 311105, got %d", records[0].MeasureValues)
	}
}

// TestProcessCSV_URLValidation tests various URL formats
func TestProcessCSV_URLValidation(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectSkipped bool
	}{
		{"valid www URL", "www.example.com/page", false},
		{"valid www with subdomain", "www.subdomain.example.com/page", false},
		{"invalid https URL", "https://example.com/page", true},
		{"invalid http URL", "http://example.com/page", true},
		{"invalid no www", "example.com/page", true},
		{"invalid relative path", "/page", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary CSV file with the test URL
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.csv")
			content := "Page,Measure Names,Measure Values\n" + tt.url + ",Pageviews,100\n"
			if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			records, err := processCSV(tmpFile)
			if err != nil {
				t.Errorf("processCSV() unexpected error for URL %q: %v", tt.url, err)
			}
			if tt.expectSkipped && len(records) != 0 {
				t.Errorf("processCSV() expected URL %q to be skipped, but got %d records", tt.url, len(records))
			}
			if !tt.expectSkipped && len(records) != 1 {
				t.Errorf("processCSV() expected URL %q to be included, but got %d records", tt.url, len(records))
			}
		})
	}
}

// TestWriteOutput_ColumnOrder tests that output has correct column order (rank, pageviews, URL)
func TestWriteOutput_ColumnOrder(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "column-order.csv")

	records := []Record{
		{Page: "www.example.com/page1", MeasureValues: 100},
	}

	err := writeOutput(records, outputPath, "1-1", 1, 1)
	if err != nil {
		t.Fatalf("writeOutput() unexpected error: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := splitLines(string(content))
	if len(lines) < 1 {
		t.Fatal("Output file is empty")
	}

	// First column should be rank (1), second should be URL
	// CSV format: "1,www.example.com/page1"
	if !contains(lines[0], "1,www.example.com/page1") {
		t.Errorf("Expected format '1,www.example.com/page1', got: %s", lines[0])
	}
}

// TestIntegration_EndToEnd tests the complete workflow with ranking
func TestIntegration_EndToEnd(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "result.csv")

	// Process the valid-with-filtering.csv file
	records, err := processCSV("testdata/valid-with-filtering.csv")
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}

	// Should get 6 Pageviews records: 300, 250, 200, 100, 50, 1 (excludes Sessions row)
	if len(records) != 6 {
		t.Fatalf("Expected 6 records, got %d", len(records))
	}

	// Write output for ranks 2-4 (should get 250, 200, 100)
	err = writeOutput(records, outputPath, "2-4", 2, 4)
	if err != nil {
		t.Fatalf("writeOutput() unexpected error: %v", err)
	}

	// Verify output file
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := splitLines(string(content))
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines in output, got %d", len(lines))
	}

	// Verify ranking (should be rank 2=page7, rank 3=page2, rank 4=page4)
	if !contains(lines[0], "2,www.example.com/page7") {
		t.Errorf("First line should be rank 2 (250 pageviews), got: %s", lines[0])
	}
	if !contains(lines[1], "3,www.example.com/page2") {
		t.Errorf("Second line should be rank 3 (200 pageviews), got: %s", lines[1])
	}
	if !contains(lines[2], "4,www.example.com/page4") {
		t.Errorf("Third line should be rank 4 (100 pageviews), got: %s", lines[2])
	}
}
