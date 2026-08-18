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
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

// dynamoDBAPI is the narrow slice of the DynamoDB client this function uses.
type dynamoDBAPI interface {
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

	log.Printf("[%s] - processing inventory reservation", ord.OrderID)

	var newInvTrans = models.Inventory{
		OrderID:    ord.OrderID,
		OrderItems: ord.ItemIds(),
	}

	// reserve the items in the inventory
	newInvTrans.Reserve()

	// Annotate saga with inventory transaction id
	ord.Inventory = newInvTrans

	// Save the reservation
	err := saveInventory(ctx, newInvTrans)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return models.Order{}, models.NewErrReserveInventory(err.Error())
	}

	// testing scenario
	if ord.OrderID[0:1] == "3" {
		return ord, models.NewErrReserveInventory("Unable to update newInvTrans for order " + ord.OrderID)
	}

	log.Printf("[%s] - reservation processed", ord.OrderID)

	return ord, nil
}

func saveInventory(ctx context.Context, newInvTrans models.Inventory) error {

	marshalledInventory, err := attributevalue.MarshalMap(newInvTrans)
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
