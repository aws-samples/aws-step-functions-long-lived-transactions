// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Inventory represents the transaction on inventory
type Inventory struct {
	TransactionID   string   `json:"transaction_id,omitempty"`
	TransactionDate string   `json:"transaction_date,omitempty"`
	OrderID         string   `json:"order_id,omitempty"`
	OrderItems      []string `json:"items,omitempty"`
	TransactionType string   `json:"transaction_type,omitempty"`
}

// Reserve method removes items from the inventory
func (i *Inventory) Reserve() {
	i.TransactionID = uuid.Must(uuid.NewV4()).String()
	i.TransactionDate = time.Now().Format(time.RFC3339)
	i.TransactionType = "Reserve"
}

// Release method makes items from the inventory available
func (i *Inventory) Release() {
	i.TransactionID = uuid.Must(uuid.NewV4()).String()
	i.TransactionDate = time.Now().Format(time.RFC3339)
	i.TransactionType = "Release"
}

/* //////////////////////////
// CUSTOM ERRORS
*/ //////////////////////////

// ErrReserveInventory represents a inventory update error
type ErrReserveInventory struct {
	message string
}

// NewErrReserveInventory constructor
func NewErrReserveInventory(message string) *ErrReserveInventory {
	return &ErrReserveInventory{
		message: message,
	}
}

func (e *ErrReserveInventory) Error() string {
	return e.message
}

// ErrReleaseInventory represents a inventory update reversal error
type ErrReleaseInventory struct {
	message string
}

// NewErrReleaseInventory constructor
func NewErrReleaseInventory(message string) *ErrReleaseInventory {
	return &ErrReleaseInventory{
		message: message,
	}
}

func (e *ErrReleaseInventory) Error() string {
	return e.message
}
