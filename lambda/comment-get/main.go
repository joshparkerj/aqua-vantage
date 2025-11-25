// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Event struct {
	Body string `json:"body"`
}

type Comment struct {
	Author string `json:"name"`
	Text   string `json:"text"`
	Date   string `json:"date"`
}

type Comments struct {
	Comments []Comment `json:"comments"`
}

var (
	client *dynamodb.Client
)

func init() {
	// TODO: Look into what the context is actually doing here 🤔
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client = dynamodb.NewFromConfig(cfg)
}

func getDynamoComment(value map[string]types.AttributeValue) Comment {
	// this should get comments from the table 😁
	date := value["Date"]
	if date != nil {
		return Comment{
			value["Author"].(*types.AttributeValueMemberS).Value,
			value["Text"].(*types.AttributeValueMemberS).Value,
			date.(*types.AttributeValueMemberS).Value,
			// I believe that if there is no "Date" in the map, then it will just come back with a zero value

		}
	} else {
		return Comment{
			value["Author"].(*types.AttributeValueMemberS).Value,
			value["Text"].(*types.AttributeValueMemberS).Value,
			"",
		}
	}
}

func getComments(ctx context.Context, tableName string) ([]Comment, error) {
	var input dynamodb.ScanInput
	var output *dynamodb.ScanOutput

	input.TableName = &tableName

	/*output*/
	// Here we scan the table of comments (to get all of them)
	output, err := client.Scan(ctx, &input)
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, output.Count)
	for i, comment := range output.Items {
		comments[i] = getDynamoComment(comment)
	}

	return comments, nil
}

func hello(ctx context.Context, event json.RawMessage) ([]byte, error) {
	var e Event

	log.Printf(string(event))

	if err := json.Unmarshal(event, &e); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return nil, err
	}

	log.Printf(e.Body)

	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		log.Printf("DYNAMO_TABLE environment variable is not set")
		return nil, fmt.Errorf("missing required environment variable DYNAMO_TABLE")
	}

	var comments []Comment

	// scan the dynamodb table using the helper method
	comments, err := getComments(ctx, tableName)
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully scanned comments from dynamodb table %s", tableName)
	// return json string as the []byte here
	jsonObj := Comments{
		comments,
	}
	jsonString, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, err
	}

	return jsonString, nil
}

func main() {
	lambda.Start(hello)
}
