package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"aws-step-functions-long-lived-transactions/models" // local

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

	marshalledOrder, err := dynamodbattribute.MarshalMap(newInvTrans)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Inventory, %v", err)
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
