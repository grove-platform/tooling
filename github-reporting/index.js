#!/usr/bin/env node

import { Command } from 'commander';
import { readFileSync } from 'fs';
import { readMetricsFromAtlas } from './read-from-db.js';
import { writeMetricsToCsv, generateOutputDir } from './write-csv.js';

const program = new Command();

program
    .name('github-reporting')
    .description('Read GitHub metrics from MongoDB Atlas and export to CSV')
    .version('1.0.0');

// Direct invocation with command-line arguments
program
    .command('export')
    .description('Export metrics to CSV using command-line arguments')
    .option('-s, --start-date <date>', 'Start date for the report (ISO format, e.g., 2024-01-01)')
    .option('-e, --end-date <date>', 'End date for the report (ISO format, e.g., 2024-12-31)')
    .option('-p, --projects <projects...>', 'List of owner/repo projects (e.g., mongodb/docs realm/realm-js)')
    .option('-o, --output <path>', 'Output directory for CSV files')
    .action(async (options) => {
        try {
            const dateRange = buildDateRange(options.startDate, options.endDate);
            const projects = parseProjects(options.projects);
            const outputDir = options.output || generateOutputDir(dateRange);

            console.log('Fetching metrics from MongoDB Atlas...');
            if (dateRange.startDate || dateRange.endDate) {
                console.log(`Date range: ${dateRange.startDate || 'beginning'} to ${dateRange.endDate || 'now'}`);
            }
            if (projects.length > 0) {
                console.log(`Projects: ${projects.map(p => `${p.owner}/${p.repo}`).join(', ')}`);
            } else {
                console.log('Projects: all');
            }

            const metrics = await readMetricsFromAtlas(dateRange, projects);
            await writeMetricsToCsv(metrics, outputDir);
        } catch (error) {
            console.error('Error:', error.message);
            process.exit(1);
        }
    });

// Config file invocation
program
    .command('export-config')
    .description('Export metrics to CSV using a configuration file')
    .argument('<config-file>', 'Path to JSON configuration file')
    .option('-o, --output <path>', 'Output directory for CSV files (overrides config file)')
    .action(async (configFile, options) => {
        try {
            const config = loadConfig(configFile);
            const dateRange = buildDateRange(config.startDate, config.endDate);
            const projects = config.projects || [];
            const outputDir = options.output || config.output || generateOutputDir(dateRange);

            console.log('Fetching metrics from MongoDB Atlas...');
            if (dateRange.startDate || dateRange.endDate) {
                console.log(`Date range: ${dateRange.startDate || 'beginning'} to ${dateRange.endDate || 'now'}`);
            }
            if (projects.length > 0) {
                console.log(`Projects: ${projects.map(p => `${p.owner}/${p.repo}`).join(', ')}`);
            } else {
                console.log('Projects: all');
            }

            const metrics = await readMetricsFromAtlas(dateRange, projects);
            await writeMetricsToCsv(metrics, outputDir);
        } catch (error) {
            console.error('Error:', error.message);
            process.exit(1);
        }
    });

/**
 * Build a date range object from start and end date strings.
 */
function buildDateRange(startDate, endDate) {
    const dateRange = {};
    if (startDate) {
        dateRange.startDate = startDate;
    }
    if (endDate) {
        dateRange.endDate = endDate;
    }
    return dateRange;
}

/**
 * Parse project strings (owner/repo format) into project objects.
 */
function parseProjects(projectStrings) {
    if (!projectStrings || projectStrings.length === 0) {
        return [];
    }

    return projectStrings.map(projectStr => {
        const parts = projectStr.split('/');
        if (parts.length !== 2) {
            throw new Error(`Invalid project format: "${projectStr}". Expected format: owner/repo`);
        }
        return { owner: parts[0], repo: parts[1] };
    });
}

/**
 * Load and parse a JSON configuration file.
 */
function loadConfig(configPath) {
    try {
        const content = readFileSync(configPath, 'utf-8');
        return JSON.parse(content);
    } catch (error) {
        if (error.code === 'ENOENT') {
            throw new Error(`Configuration file not found: ${configPath}`);
        }
        if (error instanceof SyntaxError) {
            throw new Error(`Invalid JSON in configuration file: ${error.message}`);
        }
        throw error;
    }
}

program.parse();