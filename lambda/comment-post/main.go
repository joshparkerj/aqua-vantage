// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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

type Comment struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

var (
	client      *dynamodb.Client
	emailClient *sesv2.Client
)

func init() {
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

func saveComment(ctx context.Context, tableName string, comment Comment) error {
	var input dynamodb.PutItemInput

	input.TableName = &tableName

	input.Item = make(map[string]types.AttributeValue)
	// So, the PrimaryKey field of the comment will be a uuid just like for the appointment
	input.Item["PrimaryKey"] = dynamoString(uuid.New().String())
	input.Item["Author"] = dynamoString(comment.Author)
	input.Item["Text"] = dynamoString(comment.Text)

	dateBytes, err := time.Now().MarshalJSON()
	if err != nil {
		return err
	}

	input.Item["Date"] = dynamoString(string(dateBytes))

	/*output*/
	// Here we put the item into Dynamo DB
	_, err = client.PutItem(ctx, &input)
	if err != nil {
		return err
	}

	// Here we put the item into an email
	utf8Charset := "UTF-8"

	subjectString := "Someone has made a comment!"
	subject := sesTypes.Content{
		Data:    &subjectString,
		Charset: &utf8Charset,
	}

	bodyTextString := fmt.Sprintf("Someone has made a comment at aqua-vantage.com!"+
		" Here are all the details:\n\nAuthor: %v\nText"+
		comment.Author, comment.Text)

	bodyHtmlString := fmt.Sprintf("<html><head></head><body><h1>Someone has made a comment at aqua-vantage.com!</h1>"+
		"<p>Here are all the details:</p><ol><li>Author: %v</li><li>Text: %v</li>"+
		"</body></html>", comment.Author, comment.Text)

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
	var comment Comment
	var e Event

	log.Printf(string(event))

	if err := json.Unmarshal(event, &e); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	log.Printf(e.Body)

	if err := json.Unmarshal([]byte(e.Body), &comment); err != nil {
		log.Printf("Failed to unmarshal comment: %v", err)
		return err
	}

	log.Printf("successfully unmarshalled comment from author: %v", comment.Author)

	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		log.Printf("DYNAMO_TABLE environment variable is not set")
		return fmt.Errorf("missing required environment variable DYNAMO_TABLE")
	}

	// save the comment to dynamodb using the helper method
	if err := saveComment(ctx, tableName, comment); err != nil {
		return err
	}

	log.Printf("Successfully saved comment to dynamodb table %s", tableName)
	return nil

}

func main() {
	lambda.Start(hello)
}
