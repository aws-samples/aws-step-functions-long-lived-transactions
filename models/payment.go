// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// transactionNamespace is the UUIDv5 namespace for deterministic transaction IDs.
// Deriving transaction IDs from the order (rather than minting a random ID per
// invocation) makes the DynamoDB writes idempotent: a Step Functions task retry
// overwrites the same record instead of creating a duplicate.
var transactionNamespace = uuid.Must(uuid.FromString("8b03bc51-6db8-4d24-9e5b-0f8245f5f42e"))

// TransactionID derives a deterministic ID for a transaction on an order.
func TransactionID(orderID, transactionType string) string {
	return uuid.NewV5(transactionNamespace, orderID+"|"+transactionType).String()
}

// Payment represents a customer credit card payment
type Payment struct {
	MerchantID      string  `json:"merchant_id,omitempty" dynamodbav:"merchant_id,omitempty"`
	PaymentAmount   float64 `json:"payment_amount,omitempty" dynamodbav:"payment_amount,omitempty"`
	TransactionID   string  `json:"transaction_id,omitempty" dynamodbav:"transaction_id,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty" dynamodbav:"transaction_date,omitempty"`
	OrderID         string  `json:"order_id,omitempty" dynamodbav:"order_id,omitempty"`
	PaymentType     string  `json:"payment_type,omitempty" dynamodbav:"payment_type,omitempty"`
}

// Pay customer order payment
func (p *Payment) Pay() {
	// process payment for customer order.
	// The transaction ID is deterministic (order + type) so retries are idempotent.
	p.TransactionID = TransactionID(p.OrderID, "Debit")
	p.TransactionDate = time.Now().Format(time.RFC3339)
	p.PaymentType = "Debit"

}

// Refund customer order
func (p *Payment) Refund() {
	p.TransactionID = TransactionID(p.OrderID, "Credit")
	p.TransactionDate = time.Now().Format(time.RFC3339)
	p.PaymentAmount = -(p.PaymentAmount)
	p.PaymentType = "Credit"
}

/* //////////////////////////
// CUSTOM ERRORS
*/ //////////////////////////

// ErrProcessPayment represents a process payment error
type ErrProcessPayment struct {
	message string
}

// NewErrProcessPayment constructor
func NewErrProcessPayment(message string) *ErrProcessPayment {
	return &ErrProcessPayment{
		message: message,
	}
}
func (e *ErrProcessPayment) Error() string {
	return e.message
}

// ErrProcessRefund represents a process payment refund error
type ErrProcessRefund struct {
	message string
}

// NewErrProcessRefund constructor
func NewErrProcessRefund(message string) *ErrProcessRefund {
	return &ErrProcessRefund{
		message: message,
	}
}
func (e *ErrProcessRefund) Error() string {
	return e.message
}
