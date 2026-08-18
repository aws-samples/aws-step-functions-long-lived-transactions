// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package models

import (
	"testing"
)

const testOrderID = "40063fe3-56d9-4c51-b91f-71929834ce03"

// A retried task must produce the same transaction ID so the DynamoDB write
// overwrites the prior record instead of creating a duplicate.
func TestPayIsIdempotent(t *testing.T) {
	first := Payment{OrderID: testOrderID, PaymentAmount: 56.97}
	second := Payment{OrderID: testOrderID, PaymentAmount: 56.97}

	first.Pay()
	second.Pay()

	if first.TransactionID != second.TransactionID {
		t.Errorf("Pay() must derive the same TransactionID for the same order: %s != %s", first.TransactionID, second.TransactionID)
	}
	if first.TransactionID == "" {
		t.Error("Pay() must set a TransactionID")
	}
}

func TestRefundIsIdempotent(t *testing.T) {
	first := Payment{OrderID: testOrderID, PaymentAmount: 56.97}
	second := Payment{OrderID: testOrderID, PaymentAmount: 56.97}

	first.Refund()
	second.Refund()

	if first.TransactionID != second.TransactionID {
		t.Errorf("Refund() must derive the same TransactionID for the same order: %s != %s", first.TransactionID, second.TransactionID)
	}
}

func TestDebitAndCreditIDsDiffer(t *testing.T) {
	debit := Payment{OrderID: testOrderID, PaymentAmount: 56.97}
	credit := Payment{OrderID: testOrderID, PaymentAmount: 56.97}

	debit.Pay()
	credit.Refund()

	if debit.TransactionID == credit.TransactionID {
		t.Error("debit and credit transactions for the same order must have different IDs")
	}
}

func TestReserveIsIdempotent(t *testing.T) {
	first := Inventory{OrderID: testOrderID, OrderItems: []string{"123", "234"}}
	second := Inventory{OrderID: testOrderID, OrderItems: []string{"123", "234"}}

	first.Reserve()
	second.Reserve()

	if first.TransactionID != second.TransactionID {
		t.Errorf("Reserve() must derive the same TransactionID for the same order: %s != %s", first.TransactionID, second.TransactionID)
	}
}

func TestReserveAndReleaseIDsDiffer(t *testing.T) {
	reserve := Inventory{OrderID: testOrderID}
	release := Inventory{OrderID: testOrderID}

	reserve.Reserve()
	release.Release()

	if reserve.TransactionID == release.TransactionID {
		t.Error("reserve and release transactions for the same order must have different IDs")
	}
}

func TestDifferentOrdersGetDifferentIDs(t *testing.T) {
	a := Payment{OrderID: "order-a"}
	b := Payment{OrderID: "order-b"}

	a.Pay()
	b.Pay()

	if a.TransactionID == b.TransactionID {
		t.Error("different orders must produce different transaction IDs")
	}
}
