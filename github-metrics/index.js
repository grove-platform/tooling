import { readFile } from "fs/promises";
import { getGitHubMetrics } from "./get-github-metrics.js";
import { addMetricsToAtlas } from "./write-to-db.js";
import { RepoDetails } from "./RepoDetails.js"; // Import the RepoDetails class
import { shouldRun, updateLastRun } from "./check-last-run.js";
import { sendSlackNotification } from "./send-slack-notification.js";

const dryRun = process.argv.includes("--dry-run");

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
  // Read the JSON file
  const data = await readFile("repo-details.json", "utf8");

  // Parse the JSON data into an array
  const reposArray = JSON.parse(data);

  // Convert each repo object into an instance of RepoDetails
  const repos = reposArray.map(
    (repo) => new RepoDetails(repo.owner, repo.repo),
  );

  const metricsDocs = [];

  // Iterate through the repos array
  for (const repo of repos) {
    const metricsDoc = await getGitHubMetrics(repo.owner, repo.repo);
    metricsDocs.push(metricsDoc);
  }

  if (dryRun) {
    console.log("[DRY RUN] Skipping Atlas write. Collected metrics:");
    console.log(JSON.stringify(metricsDocs, null, 2));
  } else {
    await addMetricsToAtlas(metricsDocs);

    // Update the last run timestamp after successful completion
    await updateLastRun();
  }

  return repos.length; // Return count of repos processed
}

// Main execution
async function main() {
  console.log("🚀 GitHub Metrics Collection Starting...");

  try {
    // Check if enough time has passed since last run
    if (!dryRun && !(await shouldRun())) {
      console.log("Exiting - not enough time has passed since last run");
      await sendSlackNotification({ skipped: true });
      process.exit(0);
    }

    // Process repos and collect metrics
    const repoCount = await processRepos();

    console.log("✅ GitHub Metrics Collection Complete");
    await sendSlackNotification({ success: true, repoCount });
    process.exit(0);
  } catch (error) {
    console.error("❌ Fatal error:", error);
    await sendSlackNotification({
      success: false,
      error: error.message || String(error),
    });
    process.exit(1);
  }
}

// Call the main function
main();
