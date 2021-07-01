// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Payment represents a customer credit card payment
type Payment struct {
	MerchantID      string  `json:"merchant_id,omitempty"`
	PaymentAmount   float64 `json:"payment_amount,omitempty"`
	TransactionID   string  `json:"transaction_id,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty"`
	OrderID         string  `json:"order_id,omitempty"`
	PaymentType     string  `json:"payment_type,omitempty"`
}

// Pay customer order payment
func (p *Payment) Pay() {
	// process payment for customer order
	p.TransactionID = uuid.Must(uuid.NewV4()).String()
	p.TransactionDate = time.Now().Format(time.RFC3339)
	p.PaymentType = "Debit"

}

// Refund customer order
func (p *Payment) Refund() {
	p.TransactionID = uuid.Must(uuid.NewV4()).String()
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
