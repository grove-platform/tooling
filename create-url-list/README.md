# create-url-list

A Go CLI tool that extracts and ranks URLs by pageviews from CSV data containing page analytics.

## Build

```bash
go build
```

## Usage

```bash
./create-url-list <csv-file-path> [range] [output-path]
```

### Arguments

1. **csv-file-path** (required): Path to the input CSV file
2. **range** (optional): Rank range in format `min-max` (e.g., `1-50`). Default: `1-250`
   - Specifies which ranked entries to include in the output
   - `1-50` means "get the top 50 pages by pageviews"
   - `51-100` means "get pages ranked 51-100 by pageviews"
3. **output-path** (optional): Custom output file path. Default: `output/YYYY-MM-DD_HH-MM-SS_range.csv`

### Examples

```bash
# Get top 250 pages by pageviews (default)
./create-url-list data.csv

# Get top 50 pages by pageviews
./create-url-list data.csv 1-50

# Get pages ranked 101-200 by pageviews
./create-url-list data.csv 101-200

# Specify custom output path
./create-url-list data.csv 1-100 results/top-100.csv
```

## Input Requirements

The input CSV file must contain the following columns:
- `Page`: URL of the page (must start with `www.`)
- `Measure Names`: Type of metric
- `Measure Values`: Integer value of the metric

The tool will:
- Collect all rows where `Measure Names` equals `Pageviews`
- Rank them by `Measure Values` (highest to lowest)
- Extract entries within the specified rank range
- Validate that URLs start with `www.` (to ensure consistent format without `https://`)

## Output

The output CSV file contains two columns (no headers):
1. Rank (integer) - Position in the ranking (1 = highest pageviews)
2. URL (string) - Page URL

Rows are sorted by rank in ascending order (rank 1 first).

## Error Handling

The tool exits with code 1 and displays an error message if:
- Input file path is invalid or file doesn't exist
- Required columns are missing from the CSV
- URL structure doesn't match expected format (must start with `www.`)
- Range format is invalid
