import { readFile } from 'fs/promises';
import { getGitHubMetrics } from "./get-github-metrics.js";
import { addMetricsToAtlas } from "./write-to-db.js";
import { RepoDetails } from './RepoDetails.js'; // Import the RepoDetails class
import { shouldRun, updateLastRun } from './check-last-run.js';

/* To change which repos to track metrics for, update the `repo-details.json` file.
To track metrics for a new repo, add a new entry with the owner and repo name.
You can get the owner and name from the repo URL: `https://github.com/<owner>/<repo>`
For example, to add `https://github.com/mongodb/docs-notebooks`, add:
{
  "owner": "mongodb",
  "repo": "docs-notebooks"
}
NOTE: The GitHub token used to retrieve the info from a repo MUST have repo admin permissions to access all the endpoints in this code. */

// processRepos reads the JSON config file and iterates through the repos specified, converting each to an instance of the RepoDetails class.
async function processRepos() {
    try {
        // Read the JSON file
        const data = await readFile('repo-details.json', 'utf8');

        // Parse the JSON data into an array
        const reposArray = JSON.parse(data);

        // Convert each repo object into an instance of RepoDetails
        const repos = reposArray.map(
            (repo) => new RepoDetails(repo.owner, repo.repo)
        );

        const metricsDocs = [];

        // Iterate through the repos array
        for (const repo of repos) {
            const metricsDoc = await getGitHubMetrics(repo.owner, repo.repo);
            metricsDocs.push(metricsDoc);
        }

        await addMetricsToAtlas(metricsDocs);

        // Update the last run timestamp after successful completion
        updateLastRun();
    } catch (error) {
        console.error('Error processing repos:', error);
        throw error; // Re-throw to be caught by main handler
    }
}

// Main execution
async function main() {
    console.log('🚀 GitHub Metrics Collection Starting...');

    // Check if enough time has passed since last run
    if (!shouldRun()) {
        console.log('⏭️  Exiting - not enough time has passed since last run');
        process.exit(0);
    }

    // Process repos and collect metrics
    await processRepos();

    console.log('✅ GitHub Metrics Collection Complete');
}

// Call the main function
main().catch(error => {
    console.error('❌ Fatal error:', error);
    process.exit(1);
});
