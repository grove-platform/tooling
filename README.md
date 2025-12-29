# Grove Platform Tooling

This repository contains tooling that the MongoDB Developer Docs team
uses perform various tasks related to the Grove platform.

- `audit`: Two Go projects, plus Go type definitions and constants that are
  common to both of them:
  - `gdcd`, or the Great Docs Code Devourer: an ingestion tool that gets and categorizes code examples from the current
    MongoDB documentation corpus, with a selection of metadata, and writes the info to a
    database in Atlas.
  - `dodec`, or the Database of Devoured Example Code: a query tool that lets us find code examples and related
    metadata in the database for reporting or to perform manual updates.
- `dependency-manager`: A Go CLI project to help us manage dependencies for multiple ecosystems in the docs monorepo
- `github-metrics`: a Node.js script that gets engagement metrics from GitHub for specified repos and writes them
  to a database in Atlas.
- `query-docs-feedback`: a Go project with type definitions that queries the MongoDB
  Docs Feedback received for any feedback related to code examples, and outputs the result
  to a report as `.csv`.
