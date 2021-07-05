// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"aws-step-functions-long-lived-transactions/models"

	"github.com/stretchr/testify/assert"
)

// Test Orders
var scenarioErrUpdateOrderStatus = "../testdata/scenario-2.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	t.Run("UpdateOrder", func(t *testing.T) {

		input := parseOrder(scenarioSuccessfulOrder)

		order, err := handler(nil, input)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.OrderID, "OrderID must not be empty")

	})

}
func TestErrorIsOfTypeErrProcessOrder(t *testing.T) {
	assert := assert.New(t)

	t.Run("ErrUpdateOrderStatus", func(t *testing.T) {

		o := parseOrder(scenarioErrUpdateOrderStatus)

		order, err := handler(nil, o)

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
