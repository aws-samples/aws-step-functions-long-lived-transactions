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

	log.Printf("[%s] - processing refund", ord.OrderID)

	// find Payment transaction for this order
	payment, err := getTransaction(ctx, ord.OrderID)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return ord, models.NewErrProcessRefund(err.Error())
	}

	// process the refund for the order
	payment.Refund()

	// write to database.
	err = saveTransaction(ctx, payment)
	if err != nil {
		log.Printf("[%s] - error! %s", ord.OrderID, err.Error())
		return ord, models.NewErrProcessRefund(err.Error())
	}

	// save state
	ord.Payment = payment

	// testing scenario
	if ord.OrderID[0:2] == "22" {
		return ord, models.NewErrProcessRefund("Unable to process refund for order " + ord.OrderID)
	}

	log.Printf("[%s] - refund processed", ord.OrderID)

	return ord, nil
}

func main() {
	lambda.Start(handler)
}

// getTransaction returns the debit payment transaction for the specified order
func getTransaction(ctx context.Context, orderID string) (models.Payment, error) {

	payment := models.Payment{}

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: orderID},
			":v2": &types.AttributeValueMemberS{Value: "Debit"},
		},
		KeyConditionExpression: aws.String("order_id = :v1 AND payment_type = :v2"),
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String("orderIDIndex"),
	}

	// Get payment transaction from database
	result, err := db.Query(ctx, input)
	if err != nil {
		return payment, err
	}

	if len(result.Items) == 0 {
		return payment, fmt.Errorf("no debit transaction found for order %s", orderID)
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &payment)
	if err != nil {
		return payment, fmt.Errorf("failed to DynamoDB unmarshal Payment, %w", err)
	}

	return payment, nil
}

// saveTransaction saves the refund transaction to the database
func saveTransaction(ctx context.Context, payment models.Payment) error {

	marshalledPaymentTransaction, err := attributevalue.MarshalMap(payment)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Payment, %w", err)
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledPaymentTransaction,
	})

	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %w", err)
	}
	return nil
}
