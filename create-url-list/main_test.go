package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	_, err := processCSV("testdata/missing-columns.csv", nil, "", false)
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
	records, err := processCSV("testdata/invalid-url.csv", nil, "", false)
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
			records, err := processCSV(tt.file, nil, "", false)
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
	records, err := processCSV("testdata/empty.csv", nil, "", false)
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("processCSV() got %d records, want 0", len(records))
	}
}

// TestProcessCSV_FileNotFound tests that processCSV returns an error for non-existent files
func TestProcessCSV_FileNotFound(t *testing.T) {
	_, err := processCSV("testdata/nonexistent.csv", nil, "", false)
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
				{Page: "www.example.com/page1", MeasureValues: 100, Rank: 2},
				{Page: "www.example.com/page2", MeasureValues: 50, Rank: 3},
				{Page: "www.example.com/page3", MeasureValues: 200, Rank: 1},
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
				{Page: "www.example.com/page1", MeasureValues: 100, Rank: 1},
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
			_, err := writeOutput(tt.records, tt.outputPath, tt.rangeStr, tt.minRank, tt.maxRank, false, false)
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
		{Page: "www.example.com/page3", MeasureValues: 300, Rank: 1},
		{Page: "www.example.com/page1", MeasureValues: 100, Rank: 4},
		{Page: "www.example.com/page2", MeasureValues: 200, Rank: 2},
		{Page: "www.example.com/page4", MeasureValues: 50, Rank: 3},
	}

	_, err := writeOutput(records, outputPath, "1-4", 1, 4, false, false)
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

	// Verify last data line has rank 3 (page4 with rank 3)
	if !contains(lines[3], "3,www.example.com/page4") {
		t.Errorf("Last line should contain rank 3 and page4 URL, got: %s", lines[3])
	}
}

// TestWriteOutput_NoHeaders tests that output CSV has no headers
func TestWriteOutput_NoHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "no-headers.csv")

	records := []Record{
		{Page: "www.example.com/page1", MeasureValues: 100, Rank: 1},
	}

	_, err := writeOutput(records, outputPath, "1-1", 1, 1, false, false)
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
	records, err := processCSV("testdata/more-data.csv", nil, "", false)
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

			records, err := processCSV(tmpFile, nil, "", false)
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
		{Page: "www.example.com/page1", MeasureValues: 100, Rank: 1},
	}

	_, err := writeOutput(records, outputPath, "1-1", 1, 1, false, false)
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
	records, err := processCSV("testdata/valid-with-filtering.csv", nil, "", false)
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}

	// Should get 6 Pageviews records: 300, 250, 200, 100, 50, 1 (excludes Sessions row)
	if len(records) != 6 {
		t.Fatalf("Expected 6 records, got %d", len(records))
	}

	// Write output for ranks 2-4 (should get 250, 200, 100)
	_, err = writeOutput(records, outputPath, "2-4", 2, 4, false, false)
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

// TestProcessCSV_IgnoreURLs tests that URLs in the ignore list are filtered out
func TestProcessCSV_IgnoreURLs(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.csv")
	content := `Page,Measure Names,Measure Values
www.example.com/page1,Pageviews,100
www.example.com/page2,Pageviews,200
www.example.com/page3,Pageviews,300
www.example.com/page4,Pageviews,400
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with ignore list
	ignoreURLs := []string{"www.example.com/page2", "www.example.com/page4"}
	records, err := processCSV(tmpFile, ignoreURLs, "", false)
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}

	// Should only get page1 and page3 (page2 and page4 are ignored)
	if len(records) != 2 {
		t.Errorf("processCSV() got %d records, want 2 (2 URLs ignored)", len(records))
	}

	// Verify the correct URLs are included
	foundPage1 := false
	foundPage3 := false
	for _, record := range records {
		if record.Page == "www.example.com/page1" {
			foundPage1 = true
		}
		if record.Page == "www.example.com/page3" {
			foundPage3 = true
		}
		// Make sure ignored URLs are not present
		if record.Page == "www.example.com/page2" || record.Page == "www.example.com/page4" {
			t.Errorf("Found ignored URL in results: %s", record.Page)
		}
	}

	if !foundPage1 {
		t.Error("Expected to find www.example.com/page1 in results")
	}
	if !foundPage3 {
		t.Error("Expected to find www.example.com/page3 in results")
	}
}

// TestProcessCSV_ContainsFilter tests that URLs are filtered by substring
func TestProcessCSV_ContainsFilter(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.csv")

	// Create test CSV with various URLs
	content := `Page,Measure Names,Measure Values
www.example.com/manual/page1,Pageviews,100
www.example.com/blog/post1,Pageviews,200
www.example.com/manual/page2,Pageviews,150
www.example.com/docs/guide,Pageviews,300
www.example.com/manual/tutorial,Pageviews,250
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		containsFilter string
		expectedCount  int
		expectedURLs   []string
	}{
		{
			name:           "no filter",
			containsFilter: "",
			expectedCount:  5,
			expectedURLs:   []string{"www.example.com/manual/page1", "www.example.com/blog/post1", "www.example.com/manual/page2", "www.example.com/docs/guide", "www.example.com/manual/tutorial"},
		},
		{
			name:           "filter for /manual/",
			containsFilter: "/manual/",
			expectedCount:  3,
			expectedURLs:   []string{"www.example.com/manual/page1", "www.example.com/manual/page2", "www.example.com/manual/tutorial"},
		},
		{
			name:           "filter for /blog/",
			containsFilter: "/blog/",
			expectedCount:  1,
			expectedURLs:   []string{"www.example.com/blog/post1"},
		},
		{
			name:           "filter for /docs/",
			containsFilter: "/docs/",
			expectedCount:  1,
			expectedURLs:   []string{"www.example.com/docs/guide"},
		},
		{
			name:           "filter with no matches",
			containsFilter: "/nonexistent/",
			expectedCount:  0,
			expectedURLs:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := processCSV(tmpFile, nil, tt.containsFilter, false)
			if err != nil {
				t.Fatalf("processCSV() unexpected error: %v", err)
			}

			if len(records) != tt.expectedCount {
				t.Errorf("processCSV() got %d records, want %d", len(records), tt.expectedCount)
			}

			// Verify all expected URLs are present
			for _, expectedURL := range tt.expectedURLs {
				found := false
				for _, record := range records {
					if record.Page == expectedURL {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected URL %q not found in results", expectedURL)
				}
			}

			// Verify no unexpected URLs are present
			for _, record := range records {
				found := false
				for _, expectedURL := range tt.expectedURLs {
					if record.Page == expectedURL {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Unexpected URL %q found in results", record.Page)
				}
			}
		})
	}
}

// TestProcessCSV_RankPreservation tests that original ranks are preserved after filtering
func TestProcessCSV_RankPreservation(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "rank-preservation.csv")

	// Create test CSV where filtering will skip some high-ranked items
	content := `Page,Measure Names,Measure Values
www.example.com/blog/post1,Pageviews,1000
www.example.com/blog/post2,Pageviews,900
www.example.com/manual/page1,Pageviews,800
www.example.com/blog/post3,Pageviews,700
www.example.com/manual/page2,Pageviews,600
www.example.com/blog/post4,Pageviews,500
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Filter for /manual/ URLs only
	records, err := processCSV(tmpFile, nil, "/manual/", false)
	if err != nil {
		t.Fatalf("processCSV() unexpected error: %v", err)
	}

	// Should get 2 records
	if len(records) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(records))
	}

	// Verify ranks are preserved from the original dataset
	// manual/page1 should be rank 3 (not rank 1)
	if records[0].Page == "www.example.com/manual/page1" && records[0].Rank != 3 {
		t.Errorf("manual/page1 should have rank 3, got %d", records[0].Rank)
	}

	// manual/page2 should be rank 5 (not rank 2)
	if records[1].Page == "www.example.com/manual/page2" && records[1].Rank != 5 {
		t.Errorf("manual/page2 should have rank 5, got %d", records[1].Rank)
	}

	// Verify the records are in the correct order (by rank)
	if records[0].Rank > records[1].Rank {
		t.Errorf("Records should be ordered by rank, got ranks %d and %d", records[0].Rank, records[1].Rank)
	}
}

// TestWriteOutput_ShowPageviews tests that pageviews column is added when enabled
func TestWriteOutput_ShowPageviews(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		showPageviews bool
		expectedCols  int
	}{
		{
			name:          "without pageviews",
			showPageviews: false,
			expectedCols:  2, // rank, URL
		},
		{
			name:          "with pageviews",
			showPageviews: true,
			expectedCols:  3, // rank, URL, pageviews
		},
	}

	records := []Record{
		{Page: "www.example.com/page1", MeasureValues: 300, Rank: 1},
		{Page: "www.example.com/page2", MeasureValues: 200, Rank: 2},
		{Page: "www.example.com/page3", MeasureValues: 100, Rank: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(tmpDir, tt.name+".csv")
			_, err := writeOutput(records, outputPath, "1-3", 1, 3, tt.showPageviews, false)
			if err != nil {
				t.Fatalf("writeOutput() unexpected error: %v", err)
			}

			// Read and verify the output
			content, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(string(content)), "\n")
			if len(lines) != 3 {
				t.Fatalf("Expected 3 lines, got %d", len(lines))
			}

			// Check each line has the correct number of columns
			for i, line := range lines {
				cols := strings.Split(line, ",")
				if len(cols) != tt.expectedCols {
					t.Errorf("Line %d: expected %d columns, got %d: %s", i+1, tt.expectedCols, len(cols), line)
				}

				// If showing pageviews, verify the third column is a number
				if tt.showPageviews && len(cols) == 3 {
					if _, err := strconv.Atoi(cols[2]); err != nil {
						t.Errorf("Line %d: third column should be a number, got %s", i+1, cols[2])
					}
				}
			}

			// Verify specific content when showing pageviews
			if tt.showPageviews {
				expectedLines := []string{
					"1,www.example.com/page1,300",
					"2,www.example.com/page2,200",
					"3,www.example.com/page3,100",
				}
				for i, expected := range expectedLines {
					if lines[i] != expected {
						t.Errorf("Line %d: expected %q, got %q", i+1, expected, lines[i])
					}
				}
			}
		})
	}
}

// TestWriteOutput_ShowHeaders tests that headers are added when enabled
func TestWriteOutput_ShowHeaders(t *testing.T) {
	tmpDir := t.TempDir()

	records := []Record{
		{Page: "www.example.com/page1", MeasureValues: 300, Rank: 1},
		{Page: "www.example.com/page2", MeasureValues: 200, Rank: 2},
	}

	tests := []struct {
		name            string
		showPageviews   bool
		showHeaders     bool
		expectedHeaders string
		expectedLines   int // total lines including headers
	}{
		{
			name:            "no headers, no pageviews",
			showPageviews:   false,
			showHeaders:     false,
			expectedHeaders: "",
			expectedLines:   2, // just data rows
		},
		{
			name:            "with headers, no pageviews",
			showPageviews:   false,
			showHeaders:     true,
			expectedHeaders: "Rank,Page",
			expectedLines:   3, // header + 2 data rows
		},
		{
			name:            "with headers and pageviews",
			showPageviews:   true,
			showHeaders:     true,
			expectedHeaders: "Rank,Page,Number of Page Views",
			expectedLines:   3, // header + 2 data rows
		},
		{
			name:            "no headers, with pageviews",
			showPageviews:   true,
			showHeaders:     false,
			expectedHeaders: "",
			expectedLines:   2, // just data rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(tmpDir, tt.name+".csv")
			_, err := writeOutput(records, outputPath, "1-2", 1, 2, tt.showPageviews, tt.showHeaders)
			if err != nil {
				t.Fatalf("writeOutput() unexpected error: %v", err)
			}

			// Read and verify the output
			content, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(string(content)), "\n")
			if len(lines) != tt.expectedLines {
				t.Fatalf("Expected %d lines, got %d", tt.expectedLines, len(lines))
			}

			// Check headers if they should be present
			if tt.showHeaders {
				if lines[0] != tt.expectedHeaders {
					t.Errorf("Expected headers %q, got %q", tt.expectedHeaders, lines[0])
				}
			} else {
				// First line should be data, not headers
				if !strings.HasPrefix(lines[0], "1,") {
					t.Errorf("Expected first line to start with rank '1,', got %q", lines[0])
				}
			}
		})
	}
}
