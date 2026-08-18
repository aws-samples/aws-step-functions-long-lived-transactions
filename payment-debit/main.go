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

	log.Printf("[%s] - processing payment", ord.OrderID)

	var payment = models.Payment{
		OrderID:       ord.OrderID,
		MerchantID:    "merch1",
		PaymentAmount: ord.Total(),
	}

	// Process payment
	payment.Pay()

	// Save payment
	err := savePayment(ctx, payment)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return ord, models.NewErrProcessPayment(err.Error())
	}

	// Save state
	ord.Payment = payment

	// testing scenario
	if ord.OrderID[0:1] == "2" {
		return models.Order{}, models.NewErrProcessPayment("Unable to process payment for order " + ord.OrderID)
	}

	log.Printf("[%s] - payment processed", ord.OrderID)

	return ord, nil
}

func savePayment(ctx context.Context, payment models.Payment) error {

	marshalledPayment, err := attributevalue.MarshalMap(payment)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Payment, %w", err)
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledPayment,
	})

	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %w", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
