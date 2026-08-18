// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"aws-step-functions-long-lived-transactions/models" // local

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
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
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

func handler(ctx context.Context, ord models.Order) (models.Order, error) {

	log.Printf("[%s] - processing inventory release", ord.OrderID)

	// Find inventory transaction
	inventory, err := getTransaction(ctx, ord.OrderID)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return models.Order{}, models.NewErrReleaseInventory(err.Error())
	}

	// release the items to the inventory
	inventory.Release()

	// save the inventory transaction
	err = saveTransaction(ctx, inventory)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return ord, models.NewErrReleaseInventory(err.Error())
	}

	ord.Inventory = inventory

	// testing scenario
	if ord.OrderID[0:2] == "33" {
		return ord, models.NewErrReleaseInventory("Unable to release inventory for order " + ord.OrderID)
	}

	log.Printf("[%s] - inventory release processed", ord.OrderID)

	return ord, nil
}

// getTransaction returns the reserve inventory transaction for the specified order
func getTransaction(ctx context.Context, orderID string) (models.Inventory, error) {

	inventory := models.Inventory{}

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: orderID},
			":v2": &types.AttributeValueMemberS{Value: "Reserve"},
		},
		KeyConditionExpression: aws.String("order_id = :v1 AND transaction_type = :v2"),
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String("orderIDIndex"),
	}

	// Get inventory transaction from database
	result, err := db.Query(ctx, input)
	if err != nil {
		return inventory, err
	}

	if len(result.Items) == 0 {
		return inventory, fmt.Errorf("no reserve transaction found for order %s", orderID)
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &inventory)
	if err != nil {
		return inventory, fmt.Errorf("failed to DynamoDB unmarshal Inventory, %w", err)
	}

	return inventory, nil
}

func saveTransaction(ctx context.Context, inventory models.Inventory) error {

	marshalledInventory, err := attributevalue.MarshalMap(inventory)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Inventory, %w", err)
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledInventory,
	})

	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %w", err)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
