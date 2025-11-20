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
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/google/uuid"
)

type Event struct {
	Body string `json:"body"`
}

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
	client      *dynamodb.Client
	emailClient *sesv2.Client
)

func init() {
	// Initialize the dynamodb client outside of the handler, during the init phase
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client = dynamodb.NewFromConfig(cfg)
	emailClient = sesv2.NewFromConfig(cfg)
}

func dynamoString(value string) *types.AttributeValueMemberS {
	return &types.AttributeValueMemberS{
		Value: value,
	}
}

func saveAppt(ctx context.Context, tableName string, appt Appointment) error {
	// if we're going to use a uuid, I guess it's probably at this time that we'll create it
	var input dynamodb.PutItemInput

	input.TableName = &tableName

	input.Item = make(map[string]types.AttributeValue)
	input.Item["PrimaryKey"] = dynamoString(uuid.New().String())
	input.Item["Name"] = dynamoString(appt.Name)
	input.Item["AddressLine1"] = dynamoString(appt.AddressLine1)
	input.Item["AddressLine2"] = dynamoString(appt.AddressLine2)
	input.Item["City"] = dynamoString(appt.City)
	input.Item["State"] = dynamoString(appt.State)
	input.Item["Zip"] = dynamoString(appt.Zip)
	input.Item["Phone"] = dynamoString(appt.Phone)
	input.Item["ApptDateAndTime"] = dynamoString(appt.ApptDateAndTime)

	/*output*/
	// Here we put the item into Dynamo DB
	_, err := client.PutItem(ctx, &input)
	if err != nil {
		return err
	}

	// Here we put the item into an email
	utf8Charset := "UTF-8"

	subjectString := "Someone requested an appointment!"
	subject := sesTypes.Content{
		Data:    &subjectString,
		Charset: &utf8Charset,
	}

	bodyTextString := fmt.Sprintf("Someone has requested an appointment at aqua-vantage.com!"+
		" Here are all the details:\n\nName: %v\nAddress Line 1: %v\nAddress Line 2: %v\n"+
		"City: %v\nState: %v\nZip: %v\nPhone: %v\nRequested Appointment Date and Time: %v\n"+
		"Please get into contact with the requester to confirm the appointment!", appt.Name, appt.AddressLine1, appt.AddressLine2, appt.City, appt.State, appt.Zip, appt.Phone, appt.ApptDateAndTime)

	bodyHtmlString := fmt.Sprintf("<html><head></head><body><h1>Someone has requested an appointment at aqua-vantage.com!</h1>"+
		"<p>Here are all the details:</p><ol><li>Name: %v</li><li>Address Line 1: %v</li><li>Address Line 2: %v</li>"+
		"<li>City: %v</li><li>State: %v</li><li>Zip: %v</li><li>Phone: %v</li><li>Requested Appointment Date and Time: %v</li></ol>"+
		"<p>Please get into contact with the requester to confirm the appointment!</p></body></html>", appt.Name, appt.AddressLine1, appt.AddressLine2, appt.City, appt.State, appt.Zip, appt.Phone, appt.ApptDateAndTime)

	bodyText := sesTypes.Content{
		Data:    &bodyTextString,
		Charset: &utf8Charset,
	}

	bodyHtml := sesTypes.Content{
		Data:    &bodyHtmlString,
		Charset: &utf8Charset,
	}

	body := sesTypes.Body{
		Html: &bodyHtml,
		Text: &bodyText,
	}

	message := sesTypes.Message{
		Body:    &body,
		Subject: &subject,
	}

	emailContent := sesTypes.EmailContent{
		Simple: &message,
	}

	toAddresses := make([]string, 1)
	toAddresses[0] = os.Getenv("EMAIL_RECIPIENT")

	destination := sesTypes.Destination{
		ToAddresses: toAddresses,
	}

	fromAddress := os.Getenv("EMAIL_SENDER")

	sendEmailInput := sesv2.SendEmailInput{
		Content:          &emailContent,
		Destination:      &destination,
		FromEmailAddress: &fromAddress,
	}

	sendEmailOutput, err := emailClient.SendEmail(ctx, &sendEmailInput)
	if err != nil {
		log.Printf("error sending email: %v", err)
		return err
	}

	log.Printf("Sent Email with message ID: %v", sendEmailOutput.MessageId)

	return nil
}

func hello(ctx context.Context, event json.RawMessage) error {
	var appt Appointment
	var e Event

	log.Printf(string(event))

	if err := json.Unmarshal(event, &e); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	log.Printf(e.Body)

	if err := json.Unmarshal([]byte(e.Body), &appt); err != nil {
		log.Printf("Failed to unmarshal appt: %v", err)
		return err
	}

	log.Printf("successfully unmarshalled appt for date: %v", appt.ApptDateAndTime)

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
