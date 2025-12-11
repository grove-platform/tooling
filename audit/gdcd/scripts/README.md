# Log Parser Scripts

This directory contains scripts to parse GDCD log files and analyze page changes, specifically identifying moved pages vs truly new/removed pages and tracking applied usage examples.

## Files

- `parse-log.go` - Go script that performs log parsing and analysis for page changes
- `compare-page-counts.go` - Go script that compares page counts from log files with audit-cli output
- `README.md` - This documentation file

## Purpose

### parse-log.go

The parse-log.go script analyzes log files to distinguish between:

1. **Moved Pages**: Pages that appear to be removed and created but are actually the same page moved to a new location within the same project
2. **Maybe New Pages**: Pages that may be genuinely new additions
3. **Maybe Removed Pages**: Pages that may be genuinely removed (not moved)
4. **Applied Usage Examples**: New applied usage examples on maybe new pages only

All results are reported with **project context** to clearly show which project each page belongs to.

### compare-page-counts.go

The compare-page-counts.go script compares page counts between:

1. **Log File**: Page counts extracted from GDCD log files (lines like "Found 78 docs pages for project csharp")
2. **audit-cli**: Current page counts from running `audit-cli count pages --current-only --count-by-project`

This helps identify discrepancies between what was processed during a GDCD run and the current state of the documentation repository. Differences can indicate:
- Pages added or removed since the log was generated
- Project name mismatches between systems
- Data inconsistencies that need investigation

The script automatically:
1. Runs audit-cli once to identify projects that exist only in audit-cli (not in the log)
2. Re-runs audit-cli with those projects excluded using the `--exclude-dirs` flag
3. Compares the filtered results for a cleaner comparison

The script includes built-in project name mappings to handle known differences between log file project names and audit-cli project names:
- `scala` → `scala-driver`
- `cloud-docs` → `atlas`
- `c` → `c-driver`
- `cloudgov` → `atlas-government`
- `django` → `django-mongodb`
- `docs` → `manual`
- `docs-relational-migrator` → `relational-migrator`
- `laravel` → `laravel-mongodb`
- `pymongo` → `pymongo-driver`
- `pymongo-arrow` → `pymongo-arrow-driver`
- `mck` → `kubernetes`

The script also excludes deprecated projects from comparison:
- `docs-k8s-operator` (deprecated)

## Dependencies

- Go
- `audit-cli` command (required for compare-page-counts.go) - must be available in your PATH

## How It Works

### Page Movement Detection

A page is considered "moved" if **all three conditions** are met:

1. **Same Project**: The removed page and created page are in the same project
2. **Same Code Example Count**: The removed page and created page have the same number of code examples
3. **Shared Segment**: At least one segment of the page ID (separated by `|`) is the same between the removed and created pages

For example:
- In project `ruby-driver`: `connect|tls` (removed, 6 code examples) → `security|tls` (created, 6 code examples)
- Same project AND same code examples AND shared segment `tls` → **MOVED**

### Applied Usage Examples Filtering

Applied usage examples are only counted for truly new pages, not for moved pages. This prevents double-counting when pages are reorganized.

### Maybe New and Maybe Removed Pages

Some conditions may cause moved pages to not meet our criteria for "moved" pages:

- Different number of code examples
  - Example: `connect|tls` is a "maybe removed" page and `security|tls` is a "maybe new" page but the removed page has
    6 code examples and the created page has 7 code examples
- No shared segments in page IDs
  - Example: `crud|update` is a "maybe removed" page and `write|upsert` is a "maybe new" page. Even if they have the same
    number of code examples, they share no segments in their page IDs so we can't programmatically detect that they're
    the same

Because of these conditions, we can only say that a page is "maybe new" or "maybe removed" and not "moved". A human must
manually review the "maybe new" and "maybe removed" results to determine if the page is truly new or removed. If it's
moved, we must manually adjust the count of new applied usage examples to omit the applied usage examples from the
"maybe new" but actually moved page.

## Usage

**Important**: You must be in the scripts directory to run the Go scripts directly:

### parse-log.go

```bash
# Navigate to the scripts directory first
cd /Your/Local/Filepath/tooling/audit/gdcd/scripts

# Then run the Go script
go run parse-log.go ../logs/2025-09-24-18-01-30-app.log
go run parse-log.go /absolute/path/to/your/log/file.log
```

### compare-page-counts.go

```bash
# Navigate to the scripts directory first
cd /Your/Local/Filepath/tooling/audit/gdcd/scripts

# Then run the Go script with log file and docs repo path
go run compare-page-counts.go ../logs/2025-12-10-17-58-47-app.log /path/to/docs-mongodb-internal
go run compare-page-counts.go /absolute/path/to/log/file.log /absolute/path/to/docs/repo
```

## Output Format

### parse-log.go

The parse-log.go script produces four sections:

### 1. MOVED PAGES
```
=== MOVED PAGES ===
MOVED [ruby-driver]: connect|tls -> security|tls (6 code examples)
MOVED [ruby-driver]: write|bulk-write -> crud|bulk-write (9 code examples)
MOVED [database-tools]: installation|verify -> verify (0 code examples)
```

### 2. MAYBE NEW PAGES
```
=== MAYBE NEW PAGES ===
NEW [ruby-driver]: atlas-search (2 code examples)
NEW [node]: integrations|prisma (4 code examples)
NEW [atlas-architecture]: solutions-library|rag-technical-documents (6 code examples)
```

### 3. MAYBE REMOVED PAGES
```
=== MAYBE REMOVED PAGES ===
REMOVED [ruby-driver]: common-errors (4 code examples)
REMOVED [cpp-driver]: indexes|work-with-indexes (4 code examples)
REMOVED [docs]: tutorial|install-mongodb-on-windows-unattended (11 code examples)
```

### 4. NEW APPLIED USAGE EXAMPLES
```
=== NEW APPLIED USAGE EXAMPLES ===
APPLIED USAGE [ruby-driver]: atlas-search (1 applied usage examples)
APPLIED USAGE [node]: integrations|prisma (1 applied usage examples)
APPLIED USAGE [pymongo]: data-formats|custom-types|type-codecs (1 applied usage examples)

Total new applied usage examples: 17
```

### compare-page-counts.go

The compare-page-counts.go script compares page counts from the log file with the current state from audit-cli and produces output like:

```
=== INITIAL COMPARISON ===
Found 6 projects only in audit-cli: [app-services guides mongodb-analyzer mongodb-intellij mongodb-vscode realm]

Re-running audit-cli with exclusions...

=== PAGE COUNT COMPARISON ===

Projects with differences:
--------------------------------------------------
atlas                           Log:  777    Audit:  703  (diff: -74)
atlas-architecture              Log:  124    Audit:  121  (diff: -3)
atlas-cli                       Log: 1276    Audit:  930  (diff: -346)
atlas-operator                  Log:   58    Audit:   57  (diff: -1)
c-driver                        Log:   86    Audit:   56  (diff: -30)
cloud-manager                   Log:  490    Audit:  482  (diff: -8)
compass                         Log:  117    Audit:  115  (diff: -2)
cpp-driver                      Log:   56    Audit:   52  (diff: -4)
csharp                          Log:   78    Audit:   77  (diff: -1)
database-tools                  Log:   61    Audit:   53  (diff: -8)
django-mongodb                  Log:   30    Audit:   27  (diff: -3)
drivers                         Log:   21    Audit:   20  (diff: -1)
entity-framework                Log:   13    Audit:   14  (diff: +1)
golang                          Log:  143    Audit:   68  (diff: -75)
java                            Log:   90    Audit:   89  (diff: -1)
java-rs                         Log:   56    Audit:   55  (diff: -1)
kotlin                          Log:   88    Audit:   87  (diff: -1)
kotlin-sync                     Log:   95    Audit:   66  (diff: -29)
landing                         Log:   27    Audit:   23  (diff: -4)
laravel-mongodb                 Log:   58    Audit:   57  (diff: -1)
manual                          Log: 1668    Audit: 1596  (diff: -72)
mongocli                        Log:  403    Audit:   17  (diff: -386)
mongoid                         Log:   60    Audit:   59  (diff: -1)
mongosync                       Log:   73    Audit:   88  (diff: +15)
node                            Log:   77    Audit:   76  (diff: -1)
ops-manager                     Log:  632    Audit:  628  (diff: -4)
php-library                     Log:  259    Audit:  258  (diff: -1)
pymongo-arrow-driver            Log:    8    Audit:    9  (diff: +1)
pymongo-driver                  Log:   67    Audit:   66  (diff: -1)
relational-migrator             Log:  135    Audit:  109  (diff: -26)
ruby-driver                     Log:   91    Audit:   62  (diff: -29)
rust                            Log:   76    Audit:   74  (diff: -2)
scala-driver                    Log:   44    Audit:   43  (diff: -1)
spark-connector                 Log:   16    Audit:   17  (diff: +1)
voyage                          Log:    0    Audit:    1  (diff: +1)

=== SUMMARY ===
Total projects: 43
Matching counts: 8
Different counts: 35

Total pages in log: 7869
Total pages in audit-cli: 6771
Difference: -1098
```

This helps identify:
- **Matching counts**: Projects where log and audit-cli agree
- **Different counts**: Projects where counts differ (with the difference shown)
- **Only in log**: Projects found in the log but not in audit-cli output (may indicate project name mismatches)
- **Total pages**: Sum of all page counts from each source, excluding deprecated projects and projects only in audit-cli

## Log Format Requirements

The scripts expect log lines in the following formats:

- Project context: `Project changes for <project-name>`
- Page events: `Page removed: Page ID: <page-id>` or `Page created: Page ID: <page-id>`
- Code examples: `Code example removed: Page ID: <page-id>, <count> code examples removed`
- Applied usage: `Applied usage example added: Page ID: <page-id>, <count> new applied usage examples added`

**Important**: The script tracks the current project context from "Project changes for" lines and associates all subsequent page events with that project until a new project context is encountered.
