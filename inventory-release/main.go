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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

	log.Printf("[%s] - reservation processed", ord.OrderID)

	return ord, nil
}

func getTransaction(ctx context.Context, orderID string) (models.Inventory, error) {

	inventory := models.Inventory{}

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(orderID),
			},
			":v2": {
				S: aws.String("Reserve"),
			},
		},
		KeyConditionExpression: aws.String("order_id = :v1 AND transaction_type = :v2"),
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String("orderIDIndex"),
	}

	// Get payment transaction from database
	result, err := dynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return inventory, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &inventory)
	if err != nil {
		return inventory, fmt.Errorf("failed to DynamoDB unmarshal Record, %v", err.(awserr.Error))
	}

	return inventory, nil
}

func saveTransaction(ctx context.Context, inventory models.Inventory) error {

	marshalledInventory, err := dynamodbattribute.MarshalMap(inventory)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Inventory, %v", err)
	}

	_, err = dynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledInventory,
	})

	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %v", err)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
