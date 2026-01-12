import { createObjectCsvWriter } from 'csv-writer';
import path from 'path';
import { mkdirSync } from 'fs';

/**
 * Write metrics data to three separate CSV files: summary, referrals, and top-paths.
 * @param {Array} metrics - Array of metric documents from MongoDB
 * @param {string} outputDir - Directory path for output CSV files
 * @returns {Promise<Object>} Object with paths to the created files
 */
async function writeMetricsToCsv(metrics, outputDir) {
    if (!metrics || metrics.length === 0) {
        console.log('No metrics to write.');
        return null;
    }

    // Ensure output directory exists
    mkdirSync(outputDir, { recursive: true });

    const summaryPath = path.join(outputDir, 'summary.csv');
    const referralsPath = path.join(outputDir, 'referrals.csv');
    const topPathsPath = path.join(outputDir, 'top-paths.csv');

    // Write all three files
    await Promise.all([
        writeSummaryCsv(metrics, summaryPath),
        writeReferralsCsv(metrics, referralsPath),
        writeTopPathsCsv(metrics, topPathsPath),
    ]);

    console.log(`\nSuccessfully wrote reports to ${outputDir}/`);
    console.log(`  - summary.csv (${metrics.length} records)`);

    return { summaryPath, referralsPath, topPathsPath };
}

/**
 * Write the summary CSV with core metrics (no arrays).
 */
async function writeSummaryCsv(metrics, outputPath) {
    const csvWriter = createObjectCsvWriter({
        path: outputPath,
        header: [
            { id: 'date', title: 'Date' },
            { id: 'owner', title: 'Owner' },
            { id: 'repo', title: 'Repository' },
            { id: 'clones', title: 'Clones' },
            { id: 'viewCount', title: 'Page Views' },
            { id: 'uniqueViews', title: 'Unique Views' },
            { id: 'stars', title: 'Stars' },
            { id: 'forks', title: 'Forks' },
            { id: 'watchers', title: 'Watchers' },
        ]
    });

    const records = metrics.map(metric => ({
        date: metric.date,
        owner: metric.owner,
        repo: metric.repo,
        clones: metric.clones,
        viewCount: metric.viewCount,
        uniqueViews: metric.uniqueViews,
        stars: metric.stars,
        forks: metric.forks,
        watchers: metric.watchers,
    }));

    await csvWriter.writeRecords(records);
}

/**
 * Write the referrals CSV with one row per referrer per date/repo.
 */
async function writeReferralsCsv(metrics, outputPath) {
    const csvWriter = createObjectCsvWriter({
        path: outputPath,
        header: [
            { id: 'date', title: 'Date' },
            { id: 'owner', title: 'Owner' },
            { id: 'repo', title: 'Repository' },
            { id: 'referrer', title: 'Referrer' },
            { id: 'count', title: 'Count' },
            { id: 'uniques', title: 'Uniques' },
        ]
    });

    const records = [];
    for (const metric of metrics) {
        const referralSources = metric.referralSources || [];
        for (const source of referralSources) {
            records.push({
                date: metric.date,
                owner: metric.owner,
                repo: metric.repo,
                referrer: source.referrer,
                count: source.count,
                uniques: source.uniques,
            });
        }
    }

    await csvWriter.writeRecords(records);
    console.log(`  - referrals.csv (${records.length} records)`);
}

/**
 * Write the top-paths CSV with one row per path per date/repo.
 */
async function writeTopPathsCsv(metrics, outputPath) {
    const csvWriter = createObjectCsvWriter({
        path: outputPath,
        header: [
            { id: 'date', title: 'Date' },
            { id: 'owner', title: 'Owner' },
            { id: 'repo', title: 'Repository' },
            { id: 'path', title: 'Path' },
            { id: 'count', title: 'Count' },
            { id: 'uniques', title: 'Uniques' },
        ]
    });

    const records = [];
    for (const metric of metrics) {
        const topPaths = metric.topPaths || [];
        for (const pathEntry of topPaths) {
            records.push({
                date: metric.date,
                owner: metric.owner,
                repo: metric.repo,
                path: pathEntry.path,
                count: pathEntry.count,
                uniques: pathEntry.uniques,
            });
        }
    }

    await csvWriter.writeRecords(records);
    console.log(`  - top-paths.csv (${records.length} records)`);
}

/**
 * Generate a default output directory name based on date range and timestamp.
 * @param {Object} dateRange - Object with startDate and endDate properties
 * @returns {string} Generated directory name
 */
function generateOutputDir(dateRange) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    let dirname = `github-metrics-${timestamp}`;

    if (dateRange) {
        if (dateRange.startDate) {
            const start = new Date(dateRange.startDate).toISOString().split('T')[0];
            dirname = `github-metrics-from-${start}`;
        }
        if (dateRange.endDate) {
            const end = new Date(dateRange.endDate).toISOString().split('T')[0];
            dirname += `-to-${end}`;
        }
    }

    return dirname;
}

export {
    writeMetricsToCsv,
    generateOutputDir,
}