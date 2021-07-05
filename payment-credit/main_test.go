// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/aws-samples/aws-step-functions-long-lived-transactions/models"
	"github.com/stretchr/testify/assert"
)

// Test Orders
var scenarioErrProcessRefund = "../testdata/scenario-4.json"
var scenarioSuccessfulOrder = "../testdata/scenario-7.json"

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	t.Run("ProcessRefund", func(t *testing.T) {

		input := parseOrder(scenarioSuccessfulOrder)

		order, err := handler(nil, input)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

		assert.NotEmpty(order.Payment.TransactionID, "OrderID must be empty")
	})

}

func TestErrorIsOfTypeErrProcessRefund(t *testing.T) {
	assert := assert.New(t)
	t.Run("ErrProcessRefund", func(t *testing.T) {

		input := parseOrder(scenarioErrProcessRefund)

		order, err := handler(nil, input)
		if err != nil {
			fmt.Print(err)
		}

		if assert.Error(err) {
			errorType := reflect.TypeOf(err)
			assert.Equal(errorType.String(), "*models.ErrProcessRefund", "Type does not match *models.ErrProcessRefund")
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
