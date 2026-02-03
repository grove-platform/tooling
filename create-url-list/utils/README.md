# Utility Tools for create-url-list

This directory contains diagnostic and conversion tools for troubleshooting CSV format issues with `create-url-list`.

## Tools

### debug-csv - CSV Format Inspector

Inspects a CSV file to diagnose format issues and verify column names.

**Build:**
```bash
go build -o debug-csv debug-csv.go
```

**Usage:**
```bash
./debug-csv <csv-file-path>
```

**What it shows:**
- Number of columns detected
- Exact column names (with quotes to reveal whitespace)
- Byte representation of each column name (to detect encoding issues)
- Column length
- Whether required columns (`Page`, `Measure Names`, `Measure Values`) are present
- Warnings about common issues:
  - BOM (Byte Order Mark) at start of file
  - Leading/trailing whitespace in column names

**Example:**
```bash
./debug-csv ~/Downloads/analytics-data.csv
```

**Sample output:**
```
Found 5 columns in header:

Column 0: "Page"
  Bytes: [80 97 103 101]
  Length: 4
  ✓ Matches required column 'Page'

Column 1: "Page Subsite"
  Bytes: [80 97 103 101 32 83 117 98 115 105 116 101]
  Length: 12

Column 2: "Measure Names"
  Bytes: [77 101 97 115 117 114 101 32 78 97 109 101 115]
  Length: 13
  ✓ Matches required column 'Measure Names'

Column 3: "Measure Values"
  Bytes: [77 101 97 115 117 114 101 32 86 97 108 117 101 115]
  Length: 14
  ✓ Matches required column 'Measure Values'

Required columns check:
  ✓ 'Page' found
  ✓ 'Measure Names' found
  ✓ 'Measure Values' found

✓ All required columns present!
```

---

### convert-csv - CSV Format Converter

Converts UTF-16 tab-delimited CSV files to UTF-8 comma-delimited format (standard CSV).

**Build:**
```bash
go build -o convert-csv convert-csv.go
```

**Usage:**
```bash
./convert-csv <input-file> <output-file>
```

**What it does:**
- Reads UTF-16 encoded files (with or without BOM)
- Handles tab-delimited data
- Outputs standard UTF-8 comma-delimited CSV

**Example:**
```bash
# Convert a Tableau or Excel export
./convert-csv ~/Downloads/tableau-export.csv ~/temp/converted.csv

# Then use with create-url-list
cd ..
./create-url-list ~/temp/converted.csv 1-250 output.csv
```

**Sample output:**
```
Successfully converted 51396 rows from /path/to/input.csv to /path/to/output.csv
Input format: UTF-16 tab-delimited
Output format: UTF-8 comma-delimited
```

---

## Common Workflow

When you encounter a "missing required columns" error:

1. **Diagnose the issue:**
   ```bash
   ./debug-csv /path/to/problematic-file.csv
   ```

2. **If the file is UTF-16 or tab-delimited, convert it:**
   ```bash
   ./convert-csv /path/to/problematic-file.csv /path/to/fixed-file.csv
   ```

3. **Verify the conversion worked:**
   ```bash
   ./debug-csv /path/to/fixed-file.csv
   ```

4. **Use the fixed file with create-url-list:**
   ```bash
   cd ..
   ./create-url-list /path/to/fixed-file.csv 1-250 output.csv
   ```

## Dependencies

The `convert-csv` tool requires the `golang.org/x/text` package:

```bash
go get golang.org/x/text/encoding/unicode
go get golang.org/x/text/transform
```

This dependency is automatically downloaded when you build the tool.

