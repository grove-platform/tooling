package utils

import (
	"common"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// PrintCodeExamplesByURL writes code examples found for specific URLs to a CSV file.
// Each row represents one code node with its page URL, language, file extension, category, and code length.
func PrintCodeExamplesByURL(codeExamplesMap map[string][]common.DocsPage) {
	// Generate filename with current date
	currentDate := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("wes-examples-%s.csv", currentDate)

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Page URL", "Project", "Product", "Sub-Product", "Language", "File Extension", "Category", "Code Length"}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
		return
	}

	totalPages := 0
	totalCodeNodes := 0

	// Write data rows - one row per code node
	for _, pages := range codeExamplesMap {
		for _, page := range pages {
			totalPages++
			if page.Nodes != nil {
				for _, node := range *page.Nodes {
					totalCodeNodes++
					row := []string{
						page.PageURL,
						page.ProjectName,
						page.Product,
						page.SubProduct,
						node.Language,
						node.FileExtension,
						node.Category,
						strconv.Itoa(len(node.Code)),
					}
					if err := writer.Write(row); err != nil {
						fmt.Printf("Error writing CSV row: %v\n", err)
						return
					}
				}
			}
		}
	}

	fmt.Printf("\n=== CSV Export Complete ===\n")
	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Total pages: %d\n", totalPages)
	fmt.Printf("Total code nodes: %d\n", totalCodeNodes)
}
