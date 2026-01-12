import { MongoClient } from 'mongodb';

/**
 * Read metrics from MongoDB Atlas for the given date range and projects.
 * @param {Object} dateRange - Object with startDate and endDate properties (ISO strings or Date objects)
 * @param {Array} projects - Array of {owner, repo} objects. If empty/null, reads from all collections.
 * @returns {Promise<Array>} Array of metric documents
 */
async function readMetricsFromAtlas(dateRange, projects) {
    const uri = process.env.ATLAS_CONNECTION_STRING;
    if (!uri) {
        throw new Error('ATLAS_CONNECTION_STRING environment variable is not set');
    }

    const client = new MongoClient(uri);
    let metricDocuments = [];

    try {
        await client.connect();
        const database = client.db("github_metrics");

        if (projects && projects.length > 0) {
            // Get metrics for specific projects
            for (const project of projects) {
                const projectMetrics = await getProjectMetrics(dateRange, project, database);
                metricDocuments = metricDocuments.concat(projectMetrics);
            }
        } else {
            // Get all data from all collections matching the date range
            const collections = await database.listCollections().toArray();
            for (const collInfo of collections) {
                // Skip system collections
                if (collInfo.name.startsWith('system.')) continue;

                const metrics = await getCollectionMetrics(dateRange, collInfo.name, database);
                metricDocuments = metricDocuments.concat(metrics);
            }
        }

        return metricDocuments;
    } finally {
        await client.close();
    }
}

/**
 * Get metrics for a specific project (owner/repo combination).
 * @param {Object} dateRange - Object with startDate and endDate properties
 * @param {Object} project - Object with owner and repo properties
 * @param {Db} database - MongoDB database instance
 * @returns {Promise<Array>} Array of metric documents for this project
 */
async function getProjectMetrics(dateRange, project, database) {
    const collName = project.owner + "_" + project.repo;
    return getCollectionMetrics(dateRange, collName, database);
}

/**
 * Get metrics from a specific collection within the date range.
 * @param {Object} dateRange - Object with startDate and endDate properties
 * @param {string} collName - Collection name
 * @param {Db} database - MongoDB database instance
 * @returns {Promise<Array>} Array of metric documents
 */
async function getCollectionMetrics(dateRange, collName, database) {
    try {
        const coll = database.collection(collName);

        // Build the date filter query
        const query = {};
        if (dateRange) {
            query.date = {};
            if (dateRange.startDate) {
                query.date.$gte = new Date(dateRange.startDate).toISOString();
            }
            if (dateRange.endDate) {
                query.date.$lte = new Date(dateRange.endDate).toISOString();
            }
            // If no date constraints were added, remove the empty date object
            if (Object.keys(query.date).length === 0) {
                delete query.date;
            }
        }

        const documents = await coll.find(query).sort({ date: 1 }).toArray();
        return documents;
    } catch (err) {
        console.error(`There was a problem fetching data from the '${collName}' collection: `, err);
        return [];
    }
}

export {
    readMetricsFromAtlas,
    getProjectMetrics,
    getCollectionMetrics,
}