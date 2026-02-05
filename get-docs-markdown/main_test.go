package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestExtractPageSlug tests the extractPageSlug function with various URL formats
func TestExtractPageSlug(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantSlug  string
		wantError bool
	}{
		{
			name:      "simple URL",
			url:       "www.mongodb.com/docs/manual/installation/",
			wantSlug:  "manual/installation",
			wantError: false,
		},
		{
			name:      "URL with trailing slash",
			url:       "www.mongodb.com/docs/atlas/getting-started/",
			wantSlug:  "atlas/getting-started",
			wantError: false,
		},
		{
			name:      "URL without trailing slash",
			url:       "www.mongodb.com/docs/compass/install",
			wantSlug:  "compass/install",
			wantError: false,
		},
		{
			name:      "URL with query params",
			url:       "www.mongodb.com/docs/manual/reference/?tab=cloud",
			wantSlug:  "manual/reference",
			wantError: false,
		},
		{
			name:      "URL with anchor",
			url:       "www.mongodb.com/docs/manual/reference/#section",
			wantSlug:  "manual/reference",
			wantError: false,
		},
		{
			name:      "URL with language prefix",
			url:       "www.mongodb.com/zh-cn/docs/manual/installation/",
			wantSlug:  "zh-cn/manual/installation",
			wantError: false,
		},
		{
			name:      "URL with version prefix",
			url:       "www.mongodb.com/docs/v7.0/administration/install-community/",
			wantSlug:  "v7.0/administration/install-community",
			wantError: false,
		},
		{
			name:      "URL with pt-br prefix",
			url:       "www.mongodb.com/pt-br/docs/manual/installation/",
			wantSlug:  "pt-br/manual/installation",
			wantError: false,
		},
		{
			name:      "URL with https protocol",
			url:       "https://www.mongodb.com/docs/manual/installation/",
			wantSlug:  "manual/installation",
			wantError: false,
		},
		{
			name:      "URL with http protocol",
			url:       "http://www.mongodb.com/docs/manual/installation/",
			wantSlug:  "manual/installation",
			wantError: false,
		},
		{
			name:      "URL without docs path",
			url:       "www.mongodb.com/products/",
			wantSlug:  "",
			wantError: true,
		},
		{
			name:      "URL without domain",
			url:       "/docs/manual/installation/",
			wantSlug:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := extractPageSlug(tt.url)
			if tt.wantError {
				if err == nil {
					t.Errorf("extractPageSlug(%q) expected error, got nil", tt.url)
				}
			} else {
				if err != nil {
					t.Errorf("extractPageSlug(%q) unexpected error: %v", tt.url, err)
				}
				if slug != tt.wantSlug {
					t.Errorf("extractPageSlug(%q) = %q, want %q", tt.url, slug, tt.wantSlug)
				}
			}
		})
	}
}

// TestConstructMarkdownURL tests the constructMarkdownURL function
func TestConstructMarkdownURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantURL string
	}{
		{
			name:    "simple URL",
			url:     "www.mongodb.com/docs/manual/installation/",
			wantURL: "https://www.mongodb.com/docs/manual/installation.md",
		},
		{
			name:    "URL with query params",
			url:     "www.mongodb.com/docs/manual/reference/?tab=cloud",
			wantURL: "https://www.mongodb.com/docs/manual/reference.md",
		},
		{
			name:    "URL with anchor",
			url:     "www.mongodb.com/docs/manual/reference/#section",
			wantURL: "https://www.mongodb.com/docs/manual/reference.md",
		},
		{
			name:    "URL with language prefix",
			url:     "www.mongodb.com/zh-cn/docs/manual/installation/",
			wantURL: "https://www.mongodb.com/zh-cn/docs/manual/installation.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdURL, err := constructMarkdownURL(tt.url)
			if err != nil {
				t.Errorf("constructMarkdownURL(%q) unexpected error: %v", tt.url, err)
			}
			if mdURL != tt.wantURL {
				t.Errorf("constructMarkdownURL(%q) = %q, want %q", tt.url, mdURL, tt.wantURL)
			}
		})
	}
}

// TestIsHeaderRow tests the isHeaderRow function
func TestIsHeaderRow(t *testing.T) {
	tests := []struct {
		name     string
		record   []string
		wantBool bool
	}{
		{
			name:     "output header with rank and page",
			record:   []string{"Rank", "Page"},
			wantBool: true,
		},
		{
			name:     "output header with pageviews",
			record:   []string{"Rank", "Page", "Number of Page Views"},
			wantBool: true,
		},
		{
			name:     "input header from analytics",
			record:   []string{"Page", "Page Subsite", "Measure Names", "Measure Values", "Min. Aux"},
			wantBool: true,
		},
		{
			name:     "data row with rank and URL",
			record:   []string{"1", "www.mongodb.com/docs/manual/installation/", "55197"},
			wantBool: false,
		},
		{
			name:     "data row without rank",
			record:   []string{"www.mongodb.com/docs/manual/installation/", "55197"},
			wantBool: false,
		},
		{
			name:     "header with URL keyword",
			record:   []string{"URL", "Pageviews"},
			wantBool: true,
		},
		{
			name:     "single column",
			record:   []string{"Rank"},
			wantBool: false,
		},
		{
			name:     "empty record",
			record:   []string{},
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHeaderRow(tt.record)
			if result != tt.wantBool {
				t.Errorf("isHeaderRow(%v) = %v, want %v", tt.record, result, tt.wantBool)
			}
		})
	}
}

// TestProcessCSV_FileNotFound tests that processCSV returns an error for non-existent files
func TestProcessCSV_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	// Note: processCSV requires a rate limiter, but we're just testing file opening
	// The function will fail before using the limiter
	err := processCSV("testdata/nonexistent.csv", tmpDir, 1, nil)
	if err == nil {
		t.Error("processCSV() expected error for non-existent file, got nil")
	}
}

// TestProcessCSV_EmptyFile tests that processCSV handles empty CSV files
func TestProcessCSV_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Note: processCSV requires a rate limiter, but empty file will fail before using it
	err := processCSV("testdata/empty.csv", tmpDir, 1, nil)
	// Empty file should not cause an error, just process zero URLs
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		t.Errorf("processCSV() unexpected error for empty file: %v", err)
	}
}

// Note: Full integration tests that call processCSV with actual downloads
// are skipped because they require network access and would be slow.
// The unit tests above cover the core logic without network dependencies.

// TestSlugUniqueness tests that different language/version URLs produce unique slugs
func TestSlugUniqueness(t *testing.T) {
	urls := []string{
		"www.mongodb.com/docs/manual/installation/",
		"www.mongodb.com/zh-cn/docs/manual/installation/",
		"www.mongodb.com/pt-br/docs/manual/installation/",
		"www.mongodb.com/docs/v7.0/administration/install-community/",
	}

	slugs := make(map[string]bool)
	for _, url := range urls {
		slug, err := extractPageSlug(url)
		if err != nil {
			t.Errorf("extractPageSlug(%q) unexpected error: %v", url, err)
			continue
		}

		if slugs[slug] {
			t.Errorf("Duplicate slug %q for URL %q", slug, url)
		}
		slugs[slug] = true
	}

	// Should have 4 unique slugs
	if len(slugs) != 4 {
		t.Errorf("Expected 4 unique slugs, got %d", len(slugs))
	}
}

// TestFilePathGeneration tests that file paths are generated correctly
func TestFilePathGeneration(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		outputDir    string
		expectedPath string
	}{
		{
			name:         "simple path",
			url:          "www.mongodb.com/docs/manual/installation/",
			outputDir:    "/tmp/output",
			expectedPath: "/tmp/output/manual/installation.md",
		},
		{
			name:         "with language prefix",
			url:          "www.mongodb.com/zh-cn/docs/manual/installation/",
			outputDir:    "/tmp/output",
			expectedPath: "/tmp/output/zh-cn/manual/installation.md",
		},
		{
			name:         "with version prefix",
			url:          "www.mongodb.com/docs/v7.0/administration/install/",
			outputDir:    "/tmp/output",
			expectedPath: "/tmp/output/v7.0/administration/install.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := extractPageSlug(tt.url)
			if err != nil {
				t.Fatalf("extractPageSlug(%q) unexpected error: %v", tt.url, err)
			}

			path := filepath.Join(tt.outputDir, slug+".md")
			if path != tt.expectedPath {
				t.Errorf("Expected path %q, got %q", tt.expectedPath, path)
			}
		})
	}
}
