# get-docs-markdown

A Go utility that downloads markdown versions of MongoDB documentation pages from a CSV file.

## Overview

This tool reads a CSV file containing MongoDB documentation URLs (typically output from `create-url-list`) and downloads the markdown version of each page.

## Usage

```bash
./get-docs-markdown -csv <path-to-csv> -output <output-directory> [options]
```

### Flags

- `-csv`: (Required) Path to the CSV file containing URLs
- `-output`: (Optional) Output directory for markdown files (default: `markdown-output`)
- `-workers`: (Optional) Number of concurrent download workers (default: `10`)
- `-rate-limit`: (Optional) Maximum requests per second (default: `5.0`, use `0` for unlimited)

### Examples

```bash
# Build the tool
go build

# Download markdown files with default settings (10 workers, 5 req/sec)
./get-docs-markdown -csv /path/to/top-250-dec-2025.csv -output ./markdown-files

# Use more workers and higher rate limit for faster downloads
./get-docs-markdown -csv /path/to/top-250-dec-2025.csv -output ./markdown-files -workers 20 -rate-limit 10

# Conservative settings to avoid server load
./get-docs-markdown -csv /path/to/top-250-dec-2025.csv -output ./markdown-files -workers 5 -rate-limit 2
```

## CSV Format

The tool expects a CSV file with the following format:

```
Rank,Page,Number of Page Views
1,www.mongodb.com/docs/manual/administration/install-community/,55197
2,www.mongodb.com/docs/get-started/,45669
```

Input CSVs may include or omit the header row. If a header row is present, the tool will skip it.

The tool reads the URL from the second column (index 1).

## How It Works

1. **CSV Reading**: Reads all URLs from the CSV file into memory

2. **Concurrent Processing**: Spawns multiple worker goroutines (default: 10) to download files in parallel

3. **Rate Limiting**: Uses a token bucket algorithm to limit requests per second (default: 5 req/sec)
   - Prevents overwhelming the server
   - Ensures respectful crawling behavior

4. **URL Processing**: For each URL:
   - Removes trailing slashes
   - Removes query parameters and anchor tags
   - Adds `.md` extension to get the markdown version
   - Adds User-Agent header to avoid 503 errors

5. **Slug Extraction**: Extracts the page slug from the URL (everything after `www.mongodb.com/docs/`)
   - Includes language and version prefixes to ensure uniqueness
   - Examples:
     - `www.mongodb.com/docs/manual/installation/` → `manual/installation`
     - `www.mongodb.com/zh-cn/docs/manual/installation/` → `zh-cn/manual/installation`
     - `www.mongodb.com/docs/v7.0/manual/installation/` → `v7.0/manual/installation`

6. **File Naming**: Saves files as `<output-dir>/<page-slug>.md`
   - Preserves directory structure from the URL path including language/version prefixes
   - Skips download if file already exists
   - Examples: `manual/installation.md`, `zh-cn/manual/installation.md`, `v7.0/manual/installation.md`

7. **Download**: Downloads the markdown content and saves it to the output directory

## Output

The tool creates a directory structure matching the URL paths, including language and version prefixes:

```
markdown-output/
├── manual/
│   ├── administration/
│   │   └── install-community.md
│   └── reference/
│       └── connection-string.md
├── zh-cn/
│   └── manual/
│       └── administration/
│           └── install-community.md
├── v7.0/
│   └── administration/
│       └── install-community.md
├── get-started.md
├── mongodb-shell/
│   └── install.md
└── compass/
    └── install.md
```

This ensures that different language versions and versioned documentation are saved separately without conflicts.

## Error Handling

- If a URL cannot be downloaded (404, network error, etc.), the tool logs the error and continues with the next URL
- At the end, it reports the number of successful downloads and errors

## Performance

With default settings (10 workers, 5 req/sec):
- **250 URLs**: ~50 seconds
- **500 URLs**: ~100 seconds (1.7 minutes)
- **1000 URLs**: ~200 seconds (3.3 minutes)

You can adjust `-workers` and `-rate-limit` to balance speed vs. server load. Higher values will download faster but may risk rate limiting or server errors.

