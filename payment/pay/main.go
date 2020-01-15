package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"aws-step-functions-long-lived-transactions/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/aws/aws-lambda-go/lambda"
)

var dynamoDB *dynamodb.DynamoDB

func init() {

	// Create DynamoDB client
	var awscfg = &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	var sess = session.Must(session.NewSession(awscfg))
	dynamoDB = dynamodb.New(sess)
	xray.AWS(dynamoDB.Client)

	log.SetPrefix("TRACE: ")
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

	marshalledOrder, err := dynamodbattribute.MarshalMap(payment)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Payment, %v", err)
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
