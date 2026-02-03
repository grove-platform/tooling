package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <csv-file>\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading header: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d columns in header:\n\n", len(header))
	
	for i, col := range header {
		// Show the column with quotes to reveal whitespace
		fmt.Printf("Column %d: \"%s\"\n", i, col)
		
		// Show byte representation to reveal hidden characters
		fmt.Printf("  Bytes: %v\n", []byte(col))
		
		// Show length
		fmt.Printf("  Length: %d\n", len(col))
		
		// Check for BOM or other special characters
		if strings.HasPrefix(col, "\ufeff") {
			fmt.Printf("  ⚠️  Contains BOM (Byte Order Mark) at start\n")
		}
		if strings.TrimSpace(col) != col {
			fmt.Printf("  ⚠️  Contains leading/trailing whitespace\n")
		}
		
		// Check if it matches expected columns
		switch col {
		case "Page":
			fmt.Printf("  ✓ Matches required column 'Page'\n")
		case "Measure Names":
			fmt.Printf("  ✓ Matches required column 'Measure Names'\n")
		case "Measure Values":
			fmt.Printf("  ✓ Matches required column 'Measure Values'\n")
		}
		
		fmt.Println()
	}

	// Check for required columns
	fmt.Println("Required columns check:")
	hasPage := false
	hasMeasureNames := false
	hasMeasureValues := false
	
	for _, col := range header {
		if col == "Page" {
			hasPage = true
		}
		if col == "Measure Names" {
			hasMeasureNames = true
		}
		if col == "Measure Values" {
			hasMeasureValues = true
		}
	}
	
	if hasPage {
		fmt.Println("  ✓ 'Page' found")
	} else {
		fmt.Println("  ✗ 'Page' NOT found")
	}
	
	if hasMeasureNames {
		fmt.Println("  ✓ 'Measure Names' found")
	} else {
		fmt.Println("  ✗ 'Measure Names' NOT found")
	}
	
	if hasMeasureValues {
		fmt.Println("  ✓ 'Measure Values' found")
	} else {
		fmt.Println("  ✗ 'Measure Values' NOT found")
	}
	
	if hasPage && hasMeasureNames && hasMeasureValues {
		fmt.Println("\n✓ All required columns present!")
	} else {
		fmt.Println("\n✗ Missing required columns")
	}
}

