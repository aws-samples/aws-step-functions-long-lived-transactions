package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"aws-step-functions-long-lived-transactions/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var dynamoDB *dynamodb.DynamoDB

func init() {

	// Create DynamoDB client
	var awscfg = &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	var sess = session.Must(session.NewSession(awscfg))
	dynamoDB = dynamodb.New(sess)

	// AWS X-Ray for AWS SDK trace
	xray.AWS(dynamoDB.Client)

	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime)

}

// handler for the Lambda function
func handler(ctx context.Context, ord models.Order) (models.Order, error) {

	log.Printf("[%s] - received new order", ord.OrderID)

	// persist the order data. Set order status to new
	ord.OrderStatus = "New"

	err := saveOrder(ctx, ord)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return models.Order{}, models.NewErrProcessOrder(err.Error())
	}

	// testing scenario
	if ord.OrderID[0:1] == "1" {
		return models.Order{}, models.NewErrProcessOrder("Unable to process order " + ord.OrderID)
	}

	log.Printf("[%s] - order status set to new", ord.OrderID)

	return ord, nil
}

func saveOrder(ctx context.Context, order models.Order) error {

	marshalledOrder, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Order, %v", err)
	}

	_, err = dynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledOrder,
	})
	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %v", err)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
