package aggregations

import (
	"common"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetCodeExamplesByURLs returns all usage example code nodes from documents whose page_url matches any URL in the provided array.
// The results are organized by collection name, with each collection containing an array of matching DocsPage documents.
// Only nodes with category "Usage example" are included in the results.
func GetCodeExamplesByURLs(db *mongo.Database, collectionName string, urls []string, codeExamplesMap map[string][]common.DocsPage, ctx context.Context) map[string][]common.DocsPage {
	collection := db.Collection(collectionName)

	// Expand URLs to include both www. and non-www. versions
	expandedUrls := make([]string, 0, len(urls)*2)
	for _, url := range urls {
		expandedUrls = append(expandedUrls, url)
		// Add the alternate version (with or without www.)
		if strings.Contains(url, "://www.") {
			// Has www., add version without it
			expandedUrls = append(expandedUrls, strings.Replace(url, "://www.", "://", 1))
		} else {
			// No www., add version with it
			expandedUrls = append(expandedUrls, strings.Replace(url, "://", "://www.", 1))
		}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "page_url", Value: bson.D{
				{Key: "$in", Value: expandedUrls}, // Match documents where page_url is in the array of URLs
			}},
		}}},
		// Filter nodes to only include usage examples
		{{Key: "$addFields", Value: bson.D{
			{Key: "nodes", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$nodes"},
					{Key: "as", Value: "node"},
					{Key: "cond", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$$node.category", common.UsageExample}},
					}},
				}},
			}},
		}}},
		// Update code_nodes_total to reflect filtered count
		{{Key: "$addFields", Value: bson.D{
			{Key: "code_nodes_total", Value: bson.D{
				{Key: "$size", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$nodes", bson.A{}}},
				}},
			}},
		}}},
		// Only include documents that have at least one usage example node
		{{Key: "$match", Value: bson.D{
			{Key: "code_nodes_total", Value: bson.D{
				{Key: "$gt", Value: 0},
			}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Printf("Failed to execute aggregation in collection %s: %v\n", collectionName, err)
		return codeExamplesMap
	}
	defer cursor.Close(ctx)

	var results []common.DocsPage
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Printf("Failed to decode results in collection %s: %v\n", collectionName, err)
		return codeExamplesMap
	}

	if len(results) > 0 {
		codeExamplesMap[collectionName] = results
		fmt.Printf("Found %d matching documents in %s\n", len(results), collectionName)
	}

	return codeExamplesMap
}
