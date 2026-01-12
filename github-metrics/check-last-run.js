import fs from 'fs';
import path from 'path';

// Path to the state file (mounted from persistent volume)
const STATE_FILE_PATH = process.env.STATE_FILE_PATH || '/data/last-run.json';

// Minimum days between runs
const MIN_DAYS_BETWEEN_RUNS = parseInt(process.env.MIN_DAYS_BETWEEN_RUNS || '14', 10);

/**
 * Check if enough time has passed since the last run
 * @returns {boolean} true if should run, false if should skip
 */
export function shouldRun() {
  try {
    // Check if state file exists
    if (!fs.existsSync(STATE_FILE_PATH)) {
      console.log('No previous run found. Running for the first time.');
      return true;
    }

    // Read the last run timestamp
    const stateData = JSON.parse(fs.readFileSync(STATE_FILE_PATH, 'utf8'));
    const lastRunTime = new Date(stateData.lastRun);
    const now = new Date();

    // Calculate days since last run
    const daysSinceLastRun = (now - lastRunTime) / (1000 * 60 * 60 * 24);

    console.log(`Last run: ${lastRunTime.toISOString()}`);
    console.log(`Days since last run: ${daysSinceLastRun.toFixed(2)}`);
    console.log(`Minimum days required: ${MIN_DAYS_BETWEEN_RUNS}`);

    if (daysSinceLastRun < MIN_DAYS_BETWEEN_RUNS) {
      console.log(`⏭️  Skipping run - only ${daysSinceLastRun.toFixed(2)} days since last run (need ${MIN_DAYS_BETWEEN_RUNS})`);
      return false;
    }

    console.log(`✅ Proceeding with run - ${daysSinceLastRun.toFixed(2)} days since last run`);
    return true;

  } catch (error) {
    console.error('Error checking last run time:', error.message);
    console.log('Proceeding with run due to error reading state file');
    return true; // Run if we can't read the state file
  }
}

/**
 * Update the state file with the current timestamp
 */
export function updateLastRun() {
  try {
    const now = new Date();
    const stateData = {
      lastRun: now.toISOString(),
      timestamp: now.getTime()
    };

    // Ensure the directory exists
    const dir = path.dirname(STATE_FILE_PATH);
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }

    // Write the state file
    fs.writeFileSync(STATE_FILE_PATH, JSON.stringify(stateData, null, 2), 'utf8');
    console.log(`✅ Updated last run timestamp: ${now.toISOString()}`);

  } catch (error) {
    console.error('Error updating last run time:', error.message);
    // Don't throw - we don't want to fail the job just because we can't write the state file
  }
}

/**
 * Get the last run information
 * @returns {Object|null} Object with lastRun date and timestamp, or null if no previous run
 */
export function getLastRunInfo() {
  try {
    if (!fs.existsSync(STATE_FILE_PATH)) {
      return null;
    }

    const stateData = JSON.parse(fs.readFileSync(STATE_FILE_PATH, 'utf8'));
    return {
      lastRun: new Date(stateData.lastRun),
      timestamp: stateData.timestamp
    };

  } catch (error) {
    console.error('Error reading last run info:', error.message);
    return null;
  }
}

export { MIN_DAYS_BETWEEN_RUNS };

