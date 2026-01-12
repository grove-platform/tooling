# GitHub Metrics Tooling

This directory contains tooling to enable us to track various GitHub project metrics programmatically.

This tool runs as a Kubernetes CronJob on Kanopy, automatically collecting metrics from GitHub every 14 days and storing them in MongoDB Atlas.

Planned future work:

- Add logic to work with pulled maintenance metrics once available in the test repo
- Set up Atlas Charts to visualize the data

## GitHub repo metrics

### Get metrics from GitHub

This is a simple PoC that uses [octokit](https://github.com/octokit/octokit.js) to get the following data out of GitHub
for a given repository over a trailing 14 day period:

- Views
- Unique Views
- Stars
- Watchers
- Forks
- Top 10 referral sources
- Top 10 paths/destinations in the repo

The intent is to also get the following maintenance-related stats for a given repository over a trailing 14 day period:

- Code frequency
- Commit count

However, at present, GitHub does not have any data cached for the test repo, so I'll iterate on this in a future version.

This code is in the `get-github-metrics.js` file.

> **Note**: The GitHub API does not provide the option to specify a date range for these metrics. The API _only_ provides
> this data for the trailing 14 day period, fixed. We'll need to re-run this job regularly, and in the future, we
> may want to set up a server to run this job since we cannot specify a date range.

### Change repos to track

This project pulls the configuration data from [repo-details.json](repo-details.json) to track the owner and repo name for repositories whose metrics we want to track.

#### Add a new repository

To add a new repository, create a new entry in the `repo-details.json` file in the following format:

```json
{
  "owner": "<repo-owner>",
  "repo": "<repo-name>"
}
```

You can get the owner and name from the repo URL: `https://github.com/<owner>/<repo>`
For example, to add the MongoDB docs-notebooks repository, you'd add the following to the `repo-details.json`:

```json
{
  "owner": "mongodb",
  "repo": "docs-notebooks"
}
```

The code handles tracking and writing to Atlas for multiple repos. Inserting a new repo automatically creates a new corresponding collection in Atlas. You will need to manually create Charts to correspond to the new collection.

### Write metrics to Atlas

This PoC uses the [MongoDB Node.js Driver](https://www.mongodb.com/docs/drivers/node/current/) to write the data to the
**Developer Docs** -> **Project Metrics** project in Atlas.

This code is in the `write-to-db.js` file.

In the future, we can set up Charts to visualize this data and share it with stakeholders.

### Run the tool

#### Prerequisites

To run the tool, you need:

**Atlas**:

- An Atlas Database User with write permissions for the **Developer Docs** -> **Project Metrics** project.
- A valid connection string for the cluster above.

Contact a member of the Developer Docs team to be added to this project and get the connection string.

**GitHub**:

- A [GitHub Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) (PAT) with `repo` permissions

For this project, as a MongoDB org member, you must also auth your PAT with SSO.

**System**:

- Node.js/npm installed

#### Steps

1. **Create a `.env` file**

   Create a `.env` file that contains the following details:

   ```
   ATLAS_CONNECTION_STRING="yourConnectionString"
   GITHUB_TOKEN="yourToken"
   ```

   Replace the placeholder values with your connection string and GitHub token.

   > Note: The `.env` file is in the `.gitignore`, so no worries about accidentally committing credentials.

2. **Install the dependencies**

   From the root of the directory, run the following command to install project dependencies:

   ```
   npm install
   ```

3. **Manually run the utility**

   From the root of the directory, run the following command to run the utility:

   ```
   node --env-file=.env index.js
   ```

   You should see output similar to:

   ```
   A document was inserted into mongodb_docs-notebooks with the _id: 678197a0ffe1539ff213bd86
   ```

## Automated Deployment (Kanopy CronJob)

This tool is deployed as a Kubernetes CronJob on Kanopy that runs automatically every 14 days.

### Deployment Architecture

The deployment consists of three main components:

1. **Dockerfile**: Containerizes the Node.js application
2. **cronjobs.yml**: Helm values file that configures the CronJob schedule and resources
3. **.drone.yml**: CI/CD pipeline that builds, publishes, and deploys the application

### CronJob Schedule

The cronjob is **scheduled to run weekly on Mondays at 8:00 AM UTC** (`0 8 * * 1`), but the application includes smart logic to prevent running too frequently:

- The cronjob triggers every Monday
- The application checks if 14 days have passed since the last successful run
- If less than 14 days have passed, the job exits early without collecting metrics
- If 14 days or more have passed, it collects metrics and updates the timestamp

The last run timestamp is stored in a persistent volume (`/data/last-run.json`) that survives between cronjob executions.

### Deployment Process

The deployment is fully automated via Drone CI/CD with the following steps:

1. **Check Changes**: Verifies if files in `github-metrics/` directory changed
2. **Test**: Validates dependencies with `npm ci`
3. **Build**: Builds Docker image using Kaniko and publishes to ECR
4. **Deploy**: Deploys to production Kanopy cluster using Helm
5. **Notify**: Sends Slack notification on success or failure

The pipeline only runs on pushes to the `main` branch and skips if no github-metrics files changed.

### Manual Deployment

To manually trigger a deployment:

1. Push changes to the `main` branch
2. Drone will automatically run the test, build, and deploy pipelines

### Manually Triggering the CronJob

To manually run the cronjob outside of its schedule:

```bash
# Find the cronjob
kubectl get cronjobs -n docs

# Create a one-time job from the cronjob
kubectl create job --from=cronjob/github-metrics-collection \
  github-metrics-manual-$(date +%s) -n docs

# Check the job status
kubectl get jobs -n docs

# View logs
kubectl logs -n docs job/github-metrics-manual-<timestamp>
```

### Monitoring

To check the status of the cronjob:

```bash
# View cronjob details
kubectl get cronjob github-metrics-collection -n docs

# View recent job runs
kubectl get jobs -n docs | grep github-metrics

# View logs from the most recent run
kubectl logs -n docs -l job-name=<job-name>

# Check the last run timestamp (requires exec into a pod)
kubectl exec -n docs <pod-name> -- cat /data/last-run.json
```

The logs will show whether the job ran or was skipped:
- `⏭️ Skipping run - only X days since last run (need 14)` - Job skipped, not enough time passed
- `✅ Proceeding with run - X days since last run` - Job is collecting metrics

### Configuration Changes

To modify the cronjob configuration:

1. **Change schedule**: Edit `cronjobs.yml` and update the `schedule` field
2. **Change resources**: Edit `cronjobs.yml` and update the `resources` section
3. **Change repositories tracked**: Edit `repo-details.json`

After making changes, commit and push to the `main` branch. Drone will automatically deploy the updates.
