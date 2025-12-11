package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ProjectPageCount represents the page count for a project
type ProjectPageCount struct {
	ProjectName string
	Count       int
}

// projectNameMapping maps log file project names to their audit-cli equivalents.
// This handles cases where the same project has different names in the GDCD logs
// versus the audit-cli output. Add new mappings here as needed.
var projectNameMapping = map[string]string{
	"scala":                    "scala-driver",
	"cloud-docs":               "atlas",
	"c":                        "c-driver",
	"cloudgov":                 "atlas-government",
	"django":                   "django-mongodb",
	"docs":                     "manual",
	"docs-relational-migrator": "relational-migrator",
	"laravel":                  "laravel-mongodb",
	"pymongo":                  "pymongo-driver",
	"pymongo-arrow":            "pymongo-arrow-driver",
	"mck":                      "kubernetes",
}

// deprecatedProjects lists projects that should be excluded from comparison
var deprecatedProjects = map[string]bool{
	"docs-k8s-operator": true,
}

// normalizeProjectName converts log project names to their audit-cli equivalents.
// If no mapping exists, returns the original name unchanged.
func normalizeProjectName(name string) string {
	if normalized, exists := projectNameMapping[name]; exists {
		return normalized
	}
	return name
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run compare-page-counts.go <log-file-path> <docs-repo-path>")
		fmt.Println("Example: go run compare-page-counts.go ../logs/2025-12-10-17-58-47-app.log /path/to/docs-mongodb-internal")
		os.Exit(1)
	}

	logFile := os.Args[1]
	docsRepoPath := os.Args[2]

	// Parse the log file to extract page counts
	logCounts, err := parseLogFile(logFile)
	if err != nil {
		log.Fatalf("Error parsing log file: %v", err)
	}

	// Check that audit-cli is available
	_, err = exec.LookPath("audit-cli")
	if err != nil {
		log.Fatalf("audit-cli is not available: %v", err)
	}

	// Run audit-cli command to get current page counts (first pass without exclusions)
	auditCounts, err := runAuditCli(docsRepoPath, nil)
	if err != nil {
		log.Fatalf("Error running audit-cli: %v", err)
	}

	// Find projects that are only in audit-cli
	excludeDirs := findProjectsOnlyInAudit(logCounts, auditCounts)

	// If there are projects to exclude, run audit-cli again with exclusions
	if len(excludeDirs) > 0 {
		fmt.Println("=== INITIAL COMPARISON ===")
		fmt.Printf("Found %d projects only in audit-cli: %v\n", len(excludeDirs), excludeDirs)
		fmt.Println("\nRe-running audit-cli with exclusions...\n")

		auditCounts, err = runAuditCli(docsRepoPath, excludeDirs)
		if err != nil {
			log.Fatalf("Error running audit-cli with exclusions: %v", err)
		}
	}

	// Compare the counts and report differences
	compareAndReport(logCounts, auditCounts)
}

// parseLogFile extracts page counts from the log file
func parseLogFile(logFile string) (map[string]int, error) {
	file, err := os.Open(logFile)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Regular expression to match lines like "Found 78 docs pages for project csharp"
	pageCountRegex := regexp.MustCompile(`Found (\d+) docs pages for project (.+)`)
	counts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := pageCountRegex.FindStringSubmatch(line); matches != nil {
			count, _ := strconv.Atoi(matches[1])
			projectName := strings.TrimSpace(matches[2])
			// Normalize project name to match audit-cli naming
			normalizedName := normalizeProjectName(projectName)
			// Skip deprecated projects
			if deprecatedProjects[normalizedName] {
				continue
			}
			counts[normalizedName] = count
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return counts, nil
}

// runAuditCli executes the audit-cli command and parses its output.
// If excludeDirs is provided, adds --exclude-dirs flag with comma-separated list.
func runAuditCli(docsRepoPath string, excludeDirs []string) (map[string]int, error) {
	// Build the base command
	cmdStr := fmt.Sprintf("source ~/.bashrc && audit-cli count pages %s --current-only --count-by-project", docsRepoPath)

	// Add exclude-dirs flag if provided
	if len(excludeDirs) > 0 {
		excludeList := strings.Join(excludeDirs, ",")
		cmdStr = fmt.Sprintf("%s --exclude-dirs %s", cmdStr, excludeList)
	}

	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running audit-cli: %v\nOutput: %s", err, string(output))
	}

	return parseAuditCliOutput(string(output))
}

// parseAuditCliOutput parses the output from audit-cli command
func parseAuditCliOutput(output string) (map[string]int, error) {
	counts := make(map[string]int)

	// Regular expression to match lines like "  csharp                            77"
	lineRegex := regexp.MustCompile(`^\s+([a-z0-9-]+)\s+(\d+)$`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := lineRegex.FindStringSubmatch(line); matches != nil {
			projectName := strings.TrimSpace(matches[1])
			count, _ := strconv.Atoi(matches[2])
			counts[projectName] = count
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error parsing audit-cli output: %v", err)
	}

	return counts, nil
}

// findProjectsOnlyInAudit identifies projects that exist in audit-cli but not in the log
func findProjectsOnlyInAudit(logCounts, auditCounts map[string]int) []string {
	var onlyInAudit []string

	for project := range auditCounts {
		if _, existsInLog := logCounts[project]; !existsInLog {
			onlyInAudit = append(onlyInAudit, project)
		}
	}

	// Sort for consistent output
	for i := 0; i < len(onlyInAudit); i++ {
		for j := i + 1; j < len(onlyInAudit); j++ {
			if onlyInAudit[i] > onlyInAudit[j] {
				onlyInAudit[i], onlyInAudit[j] = onlyInAudit[j], onlyInAudit[i]
			}
		}
	}

	return onlyInAudit
}

// compareAndReport compares the two sets of counts and reports differences
func compareAndReport(logCounts, auditCounts map[string]int) {
	// Collect all unique project names and sort them
	var allProjects []string
	projectSet := make(map[string]bool)
	for project := range logCounts {
		if !projectSet[project] {
			allProjects = append(allProjects, project)
			projectSet[project] = true
		}
	}
	for project := range auditCounts {
		if !projectSet[project] {
			allProjects = append(allProjects, project)
			projectSet[project] = true
		}
	}

	// Sort projects alphabetically
	for i := 0; i < len(allProjects); i++ {
		for j := i + 1; j < len(allProjects); j++ {
			if allProjects[i] > allProjects[j] {
				allProjects[i], allProjects[j] = allProjects[j], allProjects[i]
			}
		}
	}

	// Track statistics and differences
	matching := 0
	different := 0
	onlyInLog := 0
	onlyInAudit := 0
	var differences []string
	totalLogPages := 0
	totalAuditPages := 0

	// Compare counts for each project
	for _, project := range allProjects {
		logCount, inLog := logCounts[project]
		auditCount, inAudit := auditCounts[project]

		if !inLog {
			differences = append(differences, fmt.Sprintf("%-30s  Log: N/A      Audit: %4d  (only in audit-cli)", project, auditCount))
			onlyInAudit++
			totalAuditPages += auditCount
		} else if !inAudit {
			differences = append(differences, fmt.Sprintf("%-30s  Log: %4d    Audit: N/A    (only in log)", project, logCount))
			onlyInLog++
			totalLogPages += logCount
		} else if logCount != auditCount {
			diff := auditCount - logCount
			diffStr := fmt.Sprintf("%+d", diff)
			differences = append(differences, fmt.Sprintf("%-30s  Log: %4d    Audit: %4d  (diff: %s)", project, logCount, auditCount, diffStr))
			different++
			totalLogPages += logCount
			totalAuditPages += auditCount
		} else {
			matching++
			totalLogPages += logCount
			totalAuditPages += auditCount
		}
	}

	// Print results
	fmt.Println("=== PAGE COUNT COMPARISON ===\n")

	if len(differences) > 0 {
		fmt.Println("Projects with differences:")
		fmt.Println("--------------------------------------------------")
		for _, diff := range differences {
			fmt.Println(diff)
		}
	} else {
		fmt.Println("🎉 All projects have matching counts!")
	}

	// Print summary
	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("Total projects: %d\n", len(allProjects))
	fmt.Printf("Matching counts: %d\n", matching)
	fmt.Printf("Different counts: %d\n", different)
	if onlyInLog > 0 {
		fmt.Printf("Only in log: %d\n", onlyInLog)
	}
	fmt.Println()
	fmt.Printf("Total pages in log: %d\n", totalLogPages)
	fmt.Printf("Total pages in audit-cli: %d\n", totalAuditPages)
	if totalLogPages != totalAuditPages {
		diff := totalAuditPages - totalLogPages
		fmt.Printf("Difference: %+d\n", diff)
	}
}
