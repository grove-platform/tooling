package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func BackUpDb() {
	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	var dbName = os.Getenv("DB_NAME")
	var ctx = context.Background()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
	// Define the database to copy
	sourceDb := client.Database(dbName)

	// Create a db name for today's backup
	now := time.Now()
	// Format the date as "Month_Date"
	dateStr := fmt.Sprintf("%s_%d", now.Month(), now.Day())
	targetDBName := "backup_code_metrics_" + dateStr
	targetDb := client.Database(targetDBName)

	// List all collections in the source database
	collectionNames, err := sourceDb.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Error listing collections: %v", err)
	}

	// Check if a same-day or prior-day backup already exists. If so, skip the copy.
	// This prevents overwriting the current backup when re-running a job after a partial run.
	existingBackups := getBackupDbNames(client, ctx)
	if hasRecentBackup(existingBackups) {
		log.Println("A same-day or prior-day backup already exists. Skipping backup copy.")
		return
	}

	log.Println("Backing up database...")
	// Iterate over each collection, copying records from the source to target DB
	for _, collName := range collectionNames {
		sourceColl := sourceDb.Collection(collName)
		targetColl := targetDb.Collection(collName)
		// Fetch all documents from the source collection
		cursor, err := sourceColl.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalf("Error finding documents in collection %s: %v", collName, err)
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				log.Fatalf("Error closing cursor: %v", err)
			}
		}(cursor, ctx)
		var documents []interface{}
		for cursor.Next(ctx) {
			var doc bson.M
			if err = cursor.Decode(&doc); err != nil {
				log.Fatalf("Error decoding document in collection %s: %v", collName, err)
			}
			documents = append(documents, doc)
		}
		// If the source collection contains documents, insert them into the target collection
		if len(documents) > 0 {
			_, err = targetColl.InsertMany(ctx, documents)
			if err != nil {
				log.Fatalf("Error inserting documents into collection %s: %v", collName, err)
			}
			log.Printf("Copied %d documents to collection %s\n", len(documents), collName)
		}
	}
	log.Println("Successfully backed up database")

	// Drop the oldest backup only if we have at least 3 backups AND the oldest is at least 21 days
	// old. This ensures we maintain ~3 weeks of weekly history and prevents accidentally dropping
	// all backups during partial runs or tightly-clustered re-runs.
	backupNames := getBackupDbNames(client, ctx)
	const minBackupCount = 3
	const minBackupAgeDays = 21
	if len(backupNames) < minBackupCount {
		log.Printf("Only %d backup(s) exist (need %d). Skipping oldest backup drop.", len(backupNames), minBackupCount)
		return
	}
	oldestBackupDate, err := parseBackupDate(findOldestBackup(backupNames))
	if err != nil {
		log.Fatalf("Failed to parse oldest backup date: %v", err)
	}
	oldestAgedays := int(time.Since(oldestBackupDate).Hours() / 24)
	if oldestAgedays < minBackupAgeDays {
		log.Printf("Oldest backup is only %d day(s) old (need %d). Skipping oldest backup drop.", oldestAgedays, minBackupAgeDays)
		return
	}
	oldestBackupName := findOldestBackup(backupNames)
	dbToDrop := client.Database(oldestBackupName)
	err = dbToDrop.Drop(ctx)
	if err != nil {
		log.Fatalf("Failed to drop database %v: %v", oldestBackupName, err)
	}
	log.Printf("Oldest backup database '%s' dropped successfully\n", oldestBackupName)
}

// hasRecentBackup returns true if any backup is from today or yesterday.
func hasRecentBackup(backupNames []string) bool {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	for _, name := range backupNames {
		date, err := parseBackupDate(name)
		if err != nil {
			continue
		}
		if !date.Before(yesterday) {
			return true
		}
	}
	return false
}

// parseBackupDate extracts the date from a backup database name (e.g. "backup_code_metrics_March_26").
func parseBackupDate(name string) (time.Time, error) {
	parts := strings.Split(name, "_")
	if len(parts) < 4 {
		return time.Time{}, fmt.Errorf("invalid backup name: %s", name)
	}
	monthStr := parts[len(parts)-2]
	dayStr := parts[len(parts)-1]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day in backup name %s: %v", name, err)
	}
	month, err := parseMonth(monthStr)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now().UTC()
	year := now.Year()
	if month > now.Month() {
		year--
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

// The cluster contains a mix of databases - some are backups, and some are other databases.
// We want to get a slice of only the backup database names.
func getBackupDbNames(client *mongo.Client, ctx context.Context) []string {
	var backupNames []string
	// List the database names in the cluster
	databaseNames, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Failed to list database names: %v", err)
	}
	// Get only the DB names for the backup databases
	for _, databaseName := range databaseNames {
		if strings.HasPrefix(databaseName, "backup_code_metrics") {
			backupNames = append(backupNames, databaseName)
		}
	}
	return backupNames
}

// Parse the dates from the backup database names to find the oldest backup database.
func findOldestBackup(backupNames []string) string {
	var oldestDate time.Time
	var oldestBackupName string
	for _, entry := range backupNames {
		date, err := parseBackupDate(entry)
		if err != nil {
			continue
		}
		if oldestBackupName == "" || date.Before(oldestDate) {
			oldestDate = date
			oldestBackupName = entry
		}
	}
	return oldestBackupName
}

// Helper function to parse month names into time.Month
func parseMonth(month string) (time.Month, error) {
	month = strings.ToLower(month) // Make the string case-insensitive
	switch month {
	case "january":
		return time.January, nil
	case "february":
		return time.February, nil
	case "march":
		return time.March, nil
	case "april":
		return time.April, nil
	case "may":
		return time.May, nil
	case "june":
		return time.June, nil
	case "july":
		return time.July, nil
	case "august":
		return time.August, nil
	case "september":
		return time.September, nil
	case "october":
		return time.October, nil
	case "november":
		return time.November, nil
	case "december":
		return time.December, nil
	default:
		return time.Month(0), fmt.Errorf("invalid month: %s", month)
	}
}
