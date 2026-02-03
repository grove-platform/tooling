package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input-file> <output-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Converts UTF-16 tab-delimited CSV to UTF-8 comma-delimited CSV\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// Create UTF-16 decoder
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	reader := transform.NewReader(inputFile, decoder)

	// Create CSV reader with tab delimiter
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Create CSV writer (defaults to comma delimiter)
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	// Write all records
	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing record: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully converted %d rows from %s to %s\n", len(records), inputPath, outputPath)
	fmt.Printf("Input format: UTF-16 tab-delimited\n")
	fmt.Printf("Output format: UTF-8 comma-delimited\n")
}
