// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"aws-step-functions-long-lived-transactions/models"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

// Test Orders
var scenarioErrProcessPayment = "../testdata/scenario-3.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

// fakeDB satisfies dynamoDBAPI so tests run offline.
type fakeDB struct{}

func (f *fakeDB) PutItem(ctx context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	db = &fakeDB{}

	t.Run("ProcessPayment", func(t *testing.T) {

		input := parseOrder(scenarioSuccessfulOrder)

		order, err := handler(context.Background(), input)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.Payment.TransactionID, "PaymentTransactionID must not be empty")

	})
}

func TestErrorIsOfTypeErrProcessPayment(t *testing.T) {
	assert := assert.New(t)
	db = &fakeDB{}

	t.Run("ProcessPaymentErr", func(t *testing.T) {

		input := parseOrder(scenarioErrProcessPayment)

		order, err := handler(context.Background(), input)

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrProcessPayment", "Type does not match *models.ErrProcessPayment")
			assert.Empty(order.OrderID)
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
