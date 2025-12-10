# Agent Development Guide

This document provides guidance for AI agents (LLM-assisted development tools) working with the Grove Platform Tooling repository.

## Repository Overview

This is a monorepo containing multiple tools used by the MongoDB Developer Docs team for documentation-related tasks.

## Project Structure

### 1. `audit/` - Code Example Analysis Tools (Go)

Two Go projects that share common types and constants via the `audit/common` module:

#### `audit/gdcd` - Great Docs Code Devourer
- **Purpose**: Ingestion tool that extracts and categorizes code examples from MongoDB documentation
- **Language**: Go 1.24.4
- **Key Dependencies**:
  - MongoDB Go Driver v2
  - Ollama (for LLM-based code categorization using qwen2.5-coder model)
  - langchaingo
- **Module**: `module gdcd` with local replace: `replace common => ../common`
- **Build**: `go build` from `audit/gdcd/`
- **Run**: `go run .` (requires `.env` file with `MONGODB_URI` and Ollama running locally)
- **Tests**: Standard Go tests (`*_test.go` files), run with `go test ./...`
- **Long-running**: Yes (~1-2 hours depending on project count)
- **Outputs**: Logs to `logs/` directory

#### `audit/dodec` - Database of Devoured Example Code
- **Purpose**: Query tool for code example database with aggregation pipelines
- **Language**: Go 1.24.0
- **Module**: `module dodec` with local replace: `replace common => ../../common`
- **Working Directory**: `audit/dodec/src/`
- **Build**: `go build` from `audit/dodec/src/`
- **Run**: `go run .` (requires `.env` file with `MONGODB_URI`)
- **Tests**: Standard Go tests

#### `audit/common` - Shared Types
- **Purpose**: Common Go type definitions and constants
- **Module**: `module common`
- **Used by**: Both gdcd and dodec via local replace directives

### 2. `dependency-manager/` - Multi-Language Dependency Manager (Go)

- **Purpose**: CLI tool to scan and update dependencies across multiple package managers
- **Language**: Go 1.25
- **Framework**: Cobra CLI
- **Module**: `module dependency-manager`
- **Build**: `go build -o depman` from `dependency-manager/`
- **Supported Package Managers**: npm, Maven, pip, Go modules, NuGet
- **Commands**:
  - `depman check` - Dry run to check for updates
  - `depman update` - Update dependency files only
  - `depman install` - Update and install dependencies
- **Tests**: Located in `testdata/` directory
- **Documentation**: See `dependency-manager/README.md` and `dependency-manager/USAGE.md`

### 3. `github-metrics/` - GitHub Metrics Collection (Node.js)

- **Purpose**: Collects GitHub engagement metrics and writes to MongoDB Atlas
- **Language**: Node.js (ES modules)
- **Package Manager**: npm
- **Main Files**:
  - `get-github-metrics.js` - Fetches metrics from GitHub using Octokit
  - `write-to-db.js` - Writes data to MongoDB
- **Dependencies**: octokit, mongodb, esm
- **Run**: `node get-github-metrics.js` or `node write-to-db.js`
- **Status**: PoC (Proof of Concept)

### 4. `query-docs-feedback/` - Docs Feedback Query Tool (Go)

- **Purpose**: Queries MongoDB Docs Feedback for code example-related feedback
- **Language**: Go 1.23.1
- **Module**: `module query-docs-feedback`
- **Build**: `go build` from `query-docs-feedback/`
- **Run**: `go run .` (requires `.env` with `MONGODB_URI`, `DB_NAME`, `COLLECTION_NAME`)
- **Output**: CSV report

## Development Guidelines for Agents

### Go Projects

1. **Module System**: All Go projects use local module names, not GitHub paths
   - Import using local module names: `import "common"`, `import "gdcd/add-code-examples"`, etc.
   - Do NOT use full GitHub paths in imports
   - The `replace` directives in `go.mod` handle local module resolution

2. **Testing**:
   - Tests follow Go conventions: `*_test.go` files
   - Run tests with `go test ./...` from project root
   - Test data often in `test-data/` or `data/` subdirectories
   - Many projects have helper functions for testing (e.g., `GetCodeExampleForTesting()`)

3. **Environment Variables**:
   - Most projects require `.env` files (not committed to repo)
   - Common variables: `MONGODB_URI`, `DB_NAME`, `COLLECTION_NAME`
   - Use `github.com/joho/godotenv` for loading

4. **Build Commands**:
   - Always run from the project directory containing `go.mod`
   - Use `go build` or `go run .`
   - For dodec, work from `audit/dodec/src/` not `audit/dodec/`

### Node.js Projects

1. **Package Management**: Use npm (package manager commands, not manual edits)
   - Install: `npm install`
   - Add dependency: `npm install <package>`
   - Update: Use `ncu -u` then `npm install`

2. **Module System**: Uses ES modules (`"type": "module"` in package.json)

### Testing Philosophy

- Write tests for new functionality
- Run full test suite after implementation changes to catch regressions
- Remove debug output and debug files after diagnosing issues
- Optimize for maintainability over cleverness

### Code Style

- Use language-idiomatic documentation
- Capture "why" in comments, not just "what"
- Keep user-facing APIs simple (users are technical writers, not developers)
- Handle complexity internally when possible

## Common Tasks

### Running Tests
```bash
# Go projects
cd audit/gdcd && go test ./...
cd audit/dodec/src && go test ./...
cd dependency-manager && go test ./...

# Check for compilation errors
go build
```

### Building Tools
```bash
# GDCD
cd audit/gdcd && go build

# DoDEC
cd audit/dodec/src && go build

# Dependency Manager
cd dependency-manager && go build -o depman
```

### Updating Dependencies
```bash
# Go projects
go get -u ./...
go mod tidy

# Node.js projects
npm install
```

## Important Notes

- **Do NOT** manually edit `go.mod` files - use `go get` commands
- **Do NOT** manually edit `package.json` - use npm commands
- **Do NOT** create debug files without cleaning them up
- **Do NOT** add emojis or excessive success messages to output
- **Always** run full test suite after changes
- **Always** remove debug output from source code when done
