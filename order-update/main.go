// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"aws-step-functions-long-lived-transactions/models"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

// dynamoDBAPI is the narrow slice of the DynamoDB client this function uses.
type dynamoDBAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

var db dynamoDBAPI

func init() {

	// Load AWS configuration and create the DynamoDB client
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// AWS X-Ray for AWS SDK trace
	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	db = dynamodb.NewFromConfig(cfg)

	log.SetPrefix("TRACE: ")
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

	// Set order status to "pending"
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

// getOrder retrieves a specified order from DynamoDB and unmarshals it to an Order type
func getOrder(ctx context.Context, orderID string) (models.Order, error) {

	order := models.Order{}

	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"order_id": &types.AttributeValueMemberS{Value: orderID},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	result, err := db.GetItem(ctx, input)
	if err != nil {
		return order, err
	}

	err = attributevalue.UnmarshalMap(result.Item, &order)
	if err != nil {
		return order, fmt.Errorf("failed to DynamoDB unmarshal Order, %w", err)
	}

	return order, nil
}

// saveOrder persists an Order to DynamoDB
func saveOrder(ctx context.Context, order models.Order) error {

	marshalledOrder, err := attributevalue.MarshalMap(order)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Order, %w", err)
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledOrder,
	})
	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %w", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
