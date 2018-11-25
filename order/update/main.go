package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws-samples/aws-step-functions-long-lived-transactions/models"
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

	log.Printf("[%s] - received request to update order status", ord.OrderID)

	order, err := getOrder(ctx, ord.OrderID)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return ord, models.NewErrUpdateOrderStatus(err.Error())
	}

	// Set order to status to "pending"
	order.OrderStatus = "Pending"

	err = saveOrder(ctx, order)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return ord, models.NewErrUpdateOrderStatus(err.Error())
	}

	// testing scenario
	if ord.OrderID[0:2] == "11" {
		return models.Order{}, models.NewErrUpdateOrderStatus("Unable to update order status for " + ord.OrderID)
	}

	log.Printf("[%s] - order status updated to pending", ord.OrderID)

	return ord, nil
}

// getOrder retrieves a specified from DynamoDB and marshals it to a Order type
func getOrder(ctx context.Context, orderID string) (models.Order, error) {

	order := models.Order{}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"order_id": {
				S: aws.String(orderID),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	result, err := dynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return order, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &order)
	if err != nil {
		return order, fmt.Errorf("failed to DynamoDB unmarshal Order, %v", err)
	}

	return order, nil
}

// saveOrder persist an Order type to DynamoDB
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
