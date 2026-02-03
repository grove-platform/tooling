# create-url-list

A Go CLI tool that extracts and ranks URLs by pageviews from CSV data containing page analytics.

## Build

```bash
go build
```

## Usage

```bash
./create-url-list [--quiet] [--contains <substring>] <csv-file-path> [range] [output-path]
```

### Arguments

1. **--quiet** (optional): Suppress all informational output (warnings, info messages, and success messages). Only errors will be displayed. Useful when using this tool in pipelines.
2. **--contains** (optional): Filter URLs to only include those containing the specified substring. For example, `--contains /manual/` will only include URLs that contain `/manual/` in their path.
3. **csv-file-path** (required): Path to the input CSV file
4. **range** (optional): Rank range in format `min-max` (e.g., `1-50`). Default: `1-250`
   - Specifies which ranked entries to include in the output
   - `1-50` means "get the top 50 pages by pageviews"
   - `51-100` means "get pages ranked 51-100 by pageviews"
5. **output-path** (optional): Custom output file path. Default: `output/YYYY-MM-DD_HH-MM-SS_range.csv`

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

# Filter for URLs containing "/manual/" (e.g., database manual documentation)
./create-url-list --contains /manual/ data.csv

# Filter for URLs containing "/manual/" and get top 50
./create-url-list --contains /manual/ data.csv 1-50

# Use in a pipeline with quiet mode (no informational output)
./create-url-list --quiet data.csv 1-50 output.csv

# Combine multiple flags: quiet mode with URL filtering
./create-url-list --quiet --contains /manual/ data.csv 1-50 output.csv
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

## Configuration (Optional)

You can create a `config.yml` file in the same directory as the executable to configure URL filtering and output format:

```yaml
# List of URLs to ignore from the output
ignore_urls:
  - www.example.com/page-to-ignore
  - www.example.com/another-page-to-ignore

# Whether to show pageviews as a third column in the output
show_pageviews: true

# Whether to include headers in the output CSV
show_headers: true
```

### Configuration Options

**`ignore_urls`** (optional)
- URLs listed here will be completely removed from the ranking (not just hidden)
- Excluded before ranking is calculated, so remaining URLs move up without gaps
- For example, if you ignore rank #2, the former rank #3 becomes the new rank #2

**`show_pageviews`** (optional, default: `false`)
- When `false`: Output contains 2 columns (rank, URL)
- When `true`: Output contains 3 columns (rank, URL, pageviews)

**`show_headers`** (optional, default: `false`)
- When `false`: No headers in output (just data rows)
- When `true`: Adds header row with column names
  - Without pageviews: `Rank,Page`
  - With pageviews: `Rank,Page,Number of Page Views`

The config file is optional. If it doesn't exist or can't be loaded, the tool will display a warning and continue with default settings.

## Error Handling

The tool exits with code 1 and displays an error message if:
- Input file path is invalid or file doesn't exist
- Required columns are missing from the CSV
- URL structure doesn't match expected format (must start with `www.`)
- Range format is invalid

## Troubleshooting Utilities

The `utils/` directory contains helper tools for diagnosing and fixing CSV format issues.

### CSV Format Debugger

If you're getting a "missing required columns" error, use the debug tool to inspect your CSV file:

```bash
# Build the debug tool
cd utils
go build -o debug-csv debug-csv.go

# Run it on your CSV file
./debug-csv /path/to/your/file.csv
```

The debug tool will show you:
- How many columns were detected
- The exact name of each column (with quotes to reveal whitespace)
- Byte representation to reveal hidden characters (BOM, special encoding, etc.)
- Whether each required column was found
- Warnings about common issues (BOM, extra whitespace, etc.)

**Example output:**
```
Found 5 columns in header:

Column 0: "Page"
  Bytes: [80 97 103 101]
  Length: 4
  ✓ Matches required column 'Page'

Column 2: "Measure Names"
  Bytes: [77 101 97 115 117 114 101 32 78 97 109 101 115]
  Length: 13
  ✓ Matches required column 'Measure Names'
...
```

### CSV Format Converter

If your CSV file is in UTF-16 encoding or tab-delimited format (common with Excel/Tableau exports), use the converter tool:

```bash
# Build the converter tool
cd utils
go build -o convert-csv convert-csv.go

# Convert your file
./convert-csv /path/to/input.csv /path/to/output.csv
```

This tool converts:
- **From:** UTF-16 encoding with tab delimiters
- **To:** UTF-8 encoding with comma delimiters (standard CSV)

**Example:**
```bash
# Convert a Tableau export
./convert-csv ~/Downloads/tableau-export.csv ~/temp/converted.csv

# Then use the converted file
cd ..
./create-url-list ~/temp/converted.csv 1-250 output.csv
```

### Common CSV Issues

1. **UTF-16 encoding with BOM** - File starts with byte order mark (bytes `255 254`)
   - **Solution:** Use `convert-csv` tool

2. **Tab-delimited instead of comma-delimited** - Columns separated by tabs
   - **Solution:** Use `convert-csv` tool

3. **Extra whitespace in column names** - Column named `" Page "` instead of `"Page"`
   - **Solution:** Edit the CSV header row to remove extra spaces

4. **Wrong column names** - Different capitalization or spelling
   - **Solution:** Rename columns to exactly match: `Page`, `Measure Names`, `Measure Values`
