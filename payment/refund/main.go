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

	// create DynamoDB client
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

// returns a specified payment transaction from the database
func getTransaction(ctx context.Context, orderID string) (models.Payment, error) {

	payment := models.Payment{}

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(orderID),
			},
			":v2": {
				S: aws.String("Debit"),
			},
		},
		KeyConditionExpression: aws.String("order_id = :v1 AND payment_type = :v2"),
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String("orderIDIndex"),
	}

	// Get payment transaction from database
	result, err := dynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return payment, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &payment)
	if err != nil {
		return payment, fmt.Errorf("failed to DynamoDB unmarshal Payment, %v", err)
	}

	return payment, nil
}

// saves refund transaction to the database
func saveTransaction(ctx context.Context, payment models.Payment) error {

	marshalledPaymentTransaction, err := dynamodbattribute.MarshalMap(payment)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Payment, %v", err)
	}

	_, err = dynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      marshalledPaymentTransaction,
	})

	if err != nil {
		return fmt.Errorf("failed to put record to DynamoDB, %v", err)
	}
	return nil
}
