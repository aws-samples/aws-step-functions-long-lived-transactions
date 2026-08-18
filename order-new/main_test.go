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

var scenarioErrProcessOrder = "../testdata/scenario-1.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

// fakeDB satisfies dynamoDBAPI so tests run offline.
type fakeDB struct {
	putErr error
}

func (f *fakeDB) PutItem(ctx context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if f.putErr != nil {
		return nil, f.putErr
	}
	return &dynamodb.PutItemOutput{}, nil
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	db = &fakeDB{}

	t.Run("ProcessOrder", func(t *testing.T) {

		o := parseOrder(scenarioSuccessfulOrder)

		order, err := handler(context.Background(), o)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.OrderID, "OrderID must not be empty")
		assert.NotEmpty(order.CustomerID, "CustomerID must not be empty")
		assert.Equal("New", order.OrderStatus, "OrderStatus must be set to New")
		assert.True(order.Total() == 56.97, "OrderTotal does not equal expected value")
		assert.True(len(order.Items) == 3, "OrderItems should contain 3 item ids")

	})
}

func TestErrorIsOfTypeErrProcessOrder(t *testing.T) {
	assert := assert.New(t)
	db = &fakeDB{}

	t.Run("OrderProcessErr", func(t *testing.T) {

		input := parseOrder(scenarioErrProcessOrder)

		order, err := handler(context.Background(), input)

		assert.Empty(order.OrderID)

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrProcessOrder", "Type does not match *models.ErrProcessOrder")
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
