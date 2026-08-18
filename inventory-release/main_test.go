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

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// Test Orders
var scenarioErrReleaseInventory = "../testdata/scenario-4.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

// fakeDB satisfies dynamoDBAPI so tests run offline. Query returns a stored
// reserve transaction for the requested order; PutItem accepts anything.
type fakeDB struct {
	inventory models.Inventory
}

func (f *fakeDB) Query(ctx context.Context, in *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	item, err := attributevalue.MarshalMap(f.inventory)
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

	t.Run("ReleaseInventory", func(t *testing.T) {

		input := parseOrder(scenarioSuccessfulOrder)
		input.OrderID = "77063fe3-56d9-4c51-b91f-71929834ce03"
		db = &fakeDB{inventory: models.Inventory{
			OrderID:         input.OrderID,
			OrderItems:      input.ItemIds(),
			TransactionID:   "test-reserve-txn",
			TransactionType: "Reserve",
		}}

		order, err := handler(context.Background(), input)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.Inventory.TransactionID, "TransactionID must not be empty")
		assert.Equal("Release", order.Inventory.TransactionType, "TransactionType must be Release")
	})

}

func TestErrorIsOfTypeErrReleaseInventory(t *testing.T) {
	assert := assert.New(t)

	t.Run("ErrReleaseInventory", func(t *testing.T) {

		input := parseOrder(scenarioErrReleaseInventory)
		input.OrderID = "33063fe3-56d9-4c51-b91f-71929834ce03"
		db = &fakeDB{inventory: models.Inventory{
			OrderID:         input.OrderID,
			TransactionType: "Reserve",
		}}

		order, err := handler(context.Background(), input)

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrReleaseInventory", "Type does not match *models.ErrReleaseInventory")
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
