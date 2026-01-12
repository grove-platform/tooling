# GitHub Reporting

This tool reads GitHub metrics stored in MongoDB Atlas (written by the [github-metrics](../github-metrics/README.md) tool) and exports them to
CSV files for reporting to external stakeholders.

## Overview

The tool exports metrics to three separate CSV files for easy import into Google Sheets:

- **summary.csv** - Core metrics per date/repo (clones, views, stars, forks, watchers)
- **referrals.csv** - One row per referrer per date/repo
- **top-paths.csv** - One row per path per date/repo

The Date, Owner, and Repository columns in each file allow you to join/link data across sheets for analysis.

## Prerequisites

**Atlas**:

- An Atlas Database User with read permissions for the **Developer Docs** -> **Project Metrics** project.
- A valid connection string for the cluster above.

Contact a member of the Developer Docs team to be added to this project and get the connection string.

**System**:

- Node.js/npm installed

## Setup

1. **Create a `.env` file**

   Create a `.env` file that contains the following:

   ```
   ATLAS_CONNECTION_STRING="yourConnectionString"
   ```

   Replace the placeholder value with your connection string.

   > Note: The `.env` file is in the `.gitignore`, so no worries about accidentally committing credentials.

2. **Install the dependencies**

   From the root of the directory, run:

   ```
   npm install
   ```

## Usage

The tool supports two invocation methods: direct command-line arguments or a configuration file.

### Method 1: Direct Command-Line Arguments

Use the `export` command with options:

```bash
node --env-file=.env index.js export [options]
```

**Options:**

| Option | Description                                                                                      |
|--------|--------------------------------------------------------------------------------------------------|
| `-s, --start-date <date>` | Start date for the report (ISO format, e.g., 2024-01-01)                                         |
| `-e, --end-date <date>` | End date for the report (ISO format, e.g., 2024-12-31)                                           |
| `-p, --projects <projects...>` | Space-separated list of owner/repo projects (e.g., mongodb/docs mongodb/sample-app-nodejs-mflix) |
| `-o, --output <path>` | Output directory for CSV files                                                                   |

**Examples:**

```bash
# Export all metrics from all projects
node --env-file=.env index.js export -o my-report

# Export metrics for a specific date range
node --env-file=.env index.js export -s 2024-01-01 -e 2024-12-31 -o q4-report

# Export metrics for specific projects
node --env-file=.env index.js export -p mongodb/docs mongodb/docs-notebooks -o docs-report

# Combine all options
node --env-file=.env index.js export -s 2024-01-01 -e 2024-03-31 -p mongodb/docs -o q1-docs-report
```

### Method 2: Configuration File

Use the `export-config` command with a JSON configuration file:

```bash
node --env-file=.env index.js export-config <config-file> [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `-o, --output <path>` | Output directory for CSV files (overrides config file) |

**Example configuration file (`config.json`):**

```json
{
  "startDate": "2025-01-01",
  "endDate": "2025-12-31",
  "projects": [
    { "owner": "mongodb", "repo": "docs" },
    { "owner": "mongodb", "repo": "docs-notebooks" }
  ],
  "output": "annual-report"
}
```

**Run with config file:**

```bash
node --env-file=.env index.js export-config config.json
```

**Override output directory:**

```bash
node --env-file=.env index.js export-config config.json -o different-output
```

## Output

The tool creates a directory containing three CSV files:

```
my-report/
├── summary.csv
├── referrals.csv
└── top-paths.csv
```

### summary.csv

| Column | Description                                     |
|--------|-------------------------------------------------|
| Date | ISO timestamp of when metrics were collected    |
| Owner | GitHub organization/owner                       |
| Repository | Repository name                                 |
| Clones | Number of clones in the last 14 days            |
| Page Views | Total page views in the last 14 days            |
| Unique Views | Unique visitors in the last 14 days             |
| Stars | Star count (cumulative total, current count)    |
| Forks | Fork count (cumulative total, current count)    |
| Watchers | Watcher count (cumulative total, current count) |

### referrals.csv

| Column | Description |
|--------|-------------|
| Date | ISO timestamp of when metrics were collected |
| Owner | GitHub organization/owner |
| Repository | Repository name |
| Referrer | Traffic source (e.g., google.com, github.com) |
| Count | Total visits from this referrer |
| Uniques | Unique visitors from this referrer |

### top-paths.csv

| Column | Description |
|--------|-------------|
| Date | ISO timestamp of when metrics were collected |
| Owner | GitHub organization/owner |
| Repository | Repository name |
| Path | Path within the repository |
| Count | Total visits to this path |
| Uniques | Unique visitors to this path |

## Importing to Google Sheets

1. Create a new Google Sheet
2. Go to **File** → **Import**
3. Upload each CSV file as a separate sheet
4. Use the Date, Owner, and Repository columns to create relationships between sheets using VLOOKUP or pivot tables
