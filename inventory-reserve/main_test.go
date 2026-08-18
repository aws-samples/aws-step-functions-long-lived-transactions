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
var scenarioErrInventoryUpdate = "../testdata/scenario-5.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

// fakeDB satisfies dynamoDBAPI so tests run offline.
type fakeDB struct{}

func (f *fakeDB) PutItem(ctx context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	db = &fakeDB{}

	t.Run("ReserveInventory", func(t *testing.T) {

		o := parseOrder(scenarioSuccessfulOrder)

		order, err := handler(context.Background(), o)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(order.Inventory.TransactionID, "Inventory TransactionID must not be empty")

	})
}

func TestErrorIsOfTypeErrReserveInventory(t *testing.T) {
	assert := assert.New(t)
	db = &fakeDB{}

	t.Run("ReserveInventoryErr", func(t *testing.T) {

		input := parseOrder(scenarioErrInventoryUpdate)

		order, err := handler(context.Background(), input)

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrReserveInventory", "Type does not match *models.ErrReserveInventory")
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

	order := models.Order{}
	if err = jsonParser.Decode(&order); err != nil {
		println("parsing input file", err.Error())
	}

	return order
}
