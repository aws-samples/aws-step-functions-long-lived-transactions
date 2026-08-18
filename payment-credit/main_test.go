// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"aws-step-functions-long-lived-transactions/models" // local

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// Test Orders
var scenarioErrProcessRefund = "../testdata/scenario-4.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

// fakeDB satisfies dynamoDBAPI so tests run offline. Query returns a stored
// debit payment for the requested order; PutItem accepts anything.
type fakeDB struct {
	payment models.Payment
}

func (f *fakeDB) Query(ctx context.Context, in *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	item, err := attributevalue.MarshalMap(f.payment)
	if err != nil {
		return nil, err
	}
	return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
}

func (f *fakeDB) PutItem(ctx context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	t.Run("ProcessRefund", func(t *testing.T) {

		input := parseOrder(scenarioSuccessfulOrder)
		db = &fakeDB{payment: models.Payment{
			OrderID:       input.OrderID,
			MerchantID:    "merch1",
			PaymentAmount: input.Total(),
			TransactionID: "test-debit-txn",
			PaymentType:   "Debit",
		}}

		order, err := handler(context.Background(), input)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.Payment.TransactionID, "Payment TransactionID must not be empty")
		assert.Equal("Credit", order.Payment.PaymentType, "PaymentType must be Credit after refund")
	})

}

func TestErrorIsOfTypeErrProcessRefund(t *testing.T) {
	assert := assert.New(t)

	t.Run("ErrProcessRefund", func(t *testing.T) {

		input := parseOrder(scenarioErrProcessRefund)
		db = &fakeDB{payment: models.Payment{
			OrderID:     input.OrderID,
			PaymentType: "Debit",
		}}

		order, err := handler(context.Background(), input)

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrProcessRefund", "Type does not match *models.ErrProcessRefund")
			assert.NotEmpty(order.OrderID)
		}
	})
}

func parseOrder(filename string) models.Order {
	inputFile, err := os.Open(filename)
	if err != nil {
		println("opening input file", err.Error())
	}

	defer inputFile.Close()

	jsonParser := json.NewDecoder(inputFile)

	o := models.Order{}
	if err = jsonParser.Decode(&o); err != nil {
		println("parsing input file", err.Error())
	}

	return o
}
