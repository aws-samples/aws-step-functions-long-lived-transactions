// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"aws-step-functions-long-lived-transactions/models"

	"github.com/stretchr/testify/assert"
)

var scenarioErrProcessOrder = "../testdata/scenario-1.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	t.Run("ProcessOrder", func(t *testing.T) {

		o := parseOrder(scenarioSuccessfulOrder)

		order, err := handler(nil, o)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.OrderID, "OrderID must be empty")
		assert.NotEmpty(order.CustomerID, "CustomerID must not be empty")
		assert.True(order.Total() == 56.97, "OrderTotal does not equal expected value")
		assert.True(len(order.Items) == 3, "OrderItems should be contain 3 items ids")

	})
}

func TestErrorIsOfTypeErrProcessOrder(t *testing.T) {
	assert := assert.New(t)

	t.Run("OrderProcessErr", func(t *testing.T) {

		input := parseOrder(scenarioErrProcessOrder)

		order, err := handler(nil, input)
		if err != nil {
			fmt.Print(err)
		}

		assert.NotEmpty(order)

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
