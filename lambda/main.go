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
	"github.com/google/uuid"
)

type Appointment struct {
	Name            string `json:"name"`
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2"`
	City            string `json:"city"`
	State           string `json:"state"`
	Zip             string `json:"zip"`
	Phone           string `json:"phone"`
	ApptDateAndTime string `json:"apptDateAndTime"`
}

var (
	client *dynamodb.Client
)

func init() {
	// Initialize the dynamodb client outside of the handler, during the init phase
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client = dynamodb.NewFromConfig(cfg)
}

func saveAppt(ctx context.Context, tableName string, appt Appointment) error {
	// if we're going to use a uuid, I guess it's probably at this time that we'll create it
	var input dynamodb.PutItemInput

	input.TableName = &tableName

	input.Item = make(map[string]types.AttributeValue)
	input.Item["PrimaryKey"] = &types.AttributeValueMemberS{
		Value: uuid.New().String(),
	}
	input.Item["Name"] = &types.AttributeValueMemberS{
		Value: appt.Name,
	}
	input.Item["AddressLine1"] = &types.AttributeValueMemberS{
		Value: appt.AddressLine1,
	}
	input.Item["AddressLine2"] = &types.AttributeValueMemberS{
		Value: appt.AddressLine2,
	}
	input.Item["City"] = &types.AttributeValueMemberS{
		Value: appt.City,
	}
	input.Item["State"] = &types.AttributeValueMemberS{
		Value: appt.State,
	}
	input.Item["Zip"] = &types.AttributeValueMemberS{
		Value: appt.Zip,
	}
	input.Item["Phone"] = &types.AttributeValueMemberS{
		Value: appt.Phone,
	}
	input.Item["ApptDateAndTime"] = &types.AttributeValueMemberS{
		Value: appt.ApptDateAndTime,
	}

	/*output*/
	_, err := client.PutItem(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}

func hello(ctx context.Context, event json.RawMessage) error {
	var appt Appointment
	if err := json.Unmarshal(event, &appt); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		log.Printf("DYNAMO_TABLE environment variable is not set")
		return fmt.Errorf("missing required environment variable DYNAMO_TABLE")
	}

	// save the appointment to dynamodb using the helper method
	if err := saveAppt(ctx, tableName, appt); err != nil {
		return err
	}

	log.Printf("Successfully saved appt to dynamodb table %s", tableName)
	return nil

}

func main() {
	lambda.Start(hello)
}
