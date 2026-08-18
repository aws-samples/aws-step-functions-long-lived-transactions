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
	"github.com/stretchr/testify/assert"
)

// Test Orders
var scenarioErrUpdateOrderStatus = "../testdata/scenario-2.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

// fakeDB satisfies dynamoDBAPI so tests run offline. GetItem returns the
// stored order; PutItem accepts anything.
type fakeDB struct {
	stored models.Order
}

func (f *fakeDB) GetItem(ctx context.Context, in *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	item, err := attributevalue.MarshalMap(f.stored)
	if err != nil {
		return nil, err
	}
	return &dynamodb.GetItemOutput{Item: item}, nil
}

func (f *fakeDB) PutItem(ctx context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	t.Run("UpdateOrder", func(t *testing.T) {

		input := parseOrder(scenarioSuccessfulOrder)
		db = &fakeDB{stored: input}

		order, err := handler(context.Background(), input)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.OrderID, "OrderID must not be empty")

	})

}
func TestErrorIsOfTypeErrUpdateOrderStatus(t *testing.T) {
	assert := assert.New(t)

	t.Run("ErrUpdateOrderStatus", func(t *testing.T) {

		o := parseOrder(scenarioErrUpdateOrderStatus)
		db = &fakeDB{stored: o}

		order, err := handler(context.Background(), o)

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrUpdateOrderStatus", "Type does not match *models.ErrUpdateOrderStatus")
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
