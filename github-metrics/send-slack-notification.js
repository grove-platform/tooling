/**
 * Sends a Slack notification about the cronjob execution status
 * @param {Object} options - Notification options
 * @param {boolean} options.success - Whether the job succeeded
 * @param {string} [options.error] - Error message if job failed
 * @param {number} [options.repoCount] - Number of repos processed
 * @param {boolean} [options.skipped] - Whether the job was skipped due to timing
 */
export async function sendSlackNotification({ success, error, repoCount, skipped }) {
    const webhookUrl = process.env.SLACK_WEBHOOK_URL;
    
    if (!webhookUrl) {
        console.log('No SLACK_WEBHOOK_URL configured, skipping notification');
        return;
    }

    let statusEmoji, statusMsg, details;

    if (skipped) {
        statusEmoji = '⏭️';
        statusMsg = '*GitHub Metrics Collection Skipped*';
        details = 'Not enough time has passed since last run (minimum 13 days required)';
    } else if (success) {
        statusEmoji = '✅';
        statusMsg = '*GitHub Metrics Collection Succeeded*';
        details = `Successfully collected metrics for ${repoCount} repositories`;
    } else {
        statusEmoji = '❌';
        statusMsg = '*GitHub Metrics Collection Failed*';
        details = error ? `Error: ${error}` : 'Unknown error occurred';
    }

    const message = {
        text: `${statusEmoji} ${statusMsg}\n${details}\n*Time:* ${new Date().toISOString()}`
    };

    try {
        const response = await fetch(webhookUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(message),
        });

        if (!response.ok) {
            console.error(`Failed to send Slack notification: ${response.status} ${response.statusText}`);
        } else {
            console.log('Slack notification sent successfully');
        }
    } catch (err) {
        console.error('Error sending Slack notification:', err.message);
    }
}

