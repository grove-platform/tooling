import { MongoClient } from 'mongodb';

// Minimum days between runs (13 days to account for timing variations with weekly Monday runs)
// Can be overridden via MIN_DAYS_BETWEEN_RUNS environment variable
const MIN_DAYS_BETWEEN_RUNS = parseInt(process.env.MIN_DAYS_BETWEEN_RUNS || '13', 10);

// State document identifier
const STATE_DOC_ID = 'last-run-state';

/**
 * Get the MongoDB collection for storing state
 * @param {MongoClient} client - MongoDB client
 * @returns {Collection} The state collection
 */
function getStateCollection(client) {
  const database = client.db('github_metrics');
  return database.collection('job_state');
}

/**
 * Check if enough time has passed since the last run
 * @returns {Promise<boolean>} true if should run, false if should skip
 */
export async function shouldRun() {
  const uri = process.env.ATLAS_CONNECTION_STRING;
  const client = new MongoClient(uri);

  try {
    await client.connect();
    const collection = getStateCollection(client);

    // Find the last run state document
    const stateDoc = await collection.findOne({ _id: STATE_DOC_ID });

    if (!stateDoc) {
      console.log('No previous run found. Running for the first time.');
      return true;
    }

    const lastRunTime = new Date(stateDoc.lastRun);
    const now = new Date();

    // Calculate days since last run
    const daysSinceLastRun = (now - lastRunTime) / (1000 * 60 * 60 * 24);

    console.log(`Last run: ${lastRunTime.toISOString()}`);
    console.log(`Days since last run: ${daysSinceLastRun.toFixed(2)}`);
    console.log(`Minimum days required: ${MIN_DAYS_BETWEEN_RUNS}`);

    if (daysSinceLastRun < MIN_DAYS_BETWEEN_RUNS) {
      console.log(`Skipping run - only ${daysSinceLastRun.toFixed(2)} days since last run (need ${MIN_DAYS_BETWEEN_RUNS})`);
      return false;
    }

    console.log(`Proceeding with run - ${daysSinceLastRun.toFixed(2)} days since last run`);
    return true;

  } catch (error) {
    console.error('Error checking last run time:', error.message);
    console.log('Proceeding with run due to error reading state from database');
    return true; // Run if we can't read the state
  } finally {
    await client.close();
  }
}

/**
 * Update the state in MongoDB with the current timestamp
 */
export async function updateLastRun() {
  const uri = process.env.ATLAS_CONNECTION_STRING;
  const client = new MongoClient(uri);

  try {
    await client.connect();
    const collection = getStateCollection(client);
    const now = new Date();

    // Upsert the state document
    await collection.updateOne(
      { _id: STATE_DOC_ID },
      {
        $set: {
          lastRun: now.toISOString(),
          timestamp: now.getTime(),
          updatedAt: now
        }
      },
      { upsert: true }
    );

    console.log(`Updated last run timestamp in database: ${now.toISOString()}`);

  } catch (error) {
    console.error('Error updating last run time:', error.message);
    // Don't throw - we don't want to fail the job just because we can't write the state
  } finally {
    await client.close();
  }
}

export { MIN_DAYS_BETWEEN_RUNS };

