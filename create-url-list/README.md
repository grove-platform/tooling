# create-url-list

A Go CLI tool that filters CSV data containing page analytics to extract URLs with pageviews within a specified range.

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
2. **range** (optional): Integer range in format `min-max` (e.g., `1-200`). Default: `1-250`
3. **output-path** (optional): Custom output file path. Default: `output/YYYY-MM-DD_HH-MM-SS_range.csv`

### Examples

```bash
# Use default range (1-250) and default output location
./create-url-list data.csv

# Specify custom range
./create-url-list data.csv 1-500

# Specify custom range and output path
./create-url-list data.csv 100-1000 results/filtered.csv
```

## Input Requirements

The input CSV file must contain the following columns:
- `Page`: URL of the page (must start with `www.`)
- `Measure Names`: Type of metric
- `Measure Values`: Integer value of the metric

The tool will:
- Filter rows where `Measure Names` equals `Pageviews`
- Filter rows where `Measure Values` falls within the specified range (inclusive)
- Validate that URLs start with `www.` (to ensure consistent format without `https://`)

## Output

The output CSV file contains two columns (no headers):
1. `Measure Values` (integer)
2. `Page` (URL)

Rows are sorted by `Measure Values` in ascending order (lowest to highest).

## Error Handling

The tool exits with code 1 and displays an error message if:
- Input file path is invalid or file doesn't exist
- Required columns are missing from the CSV
- URL structure doesn't match expected format (must start with `www.`)
- Range format is invalid

