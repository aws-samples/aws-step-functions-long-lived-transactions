package models

import (
	"time"
)

// Order represents a customer order
type Order struct {
	OrderID     string    `json:"order_id,omitempty"`
	OrderDate   time.Time `json:"order_date,omitempty"`
	CustomerID  string    `json:"customer_id,omitempty"`
	OrderStatus string    `json:"order_status,omitempty"`
	Items       []Item    `json:"items,omitempty"`
	Payment     Payment   `json:"payment,omitempty"`
	Inventory   Inventory `json:"inventory,omitempty"`
}

// Item represents an item in the cart
type Item struct {
	ItemID      string  `json:"item_id,omitempty"`
	Qty         float64 `json:"qty,omitempty"`
	Description string  `json:"description,omitempty"`
	UnitPrice   float64 `json:"unit_price,omitempty"`
}

// Total returns the total ammount of the order
func (o Order) Total() float64 {
	var total = 0.0
	for i := 0; i <= len(o.Items)-1; i++ {

		item := o.Items[i]
		total += item.UnitPrice * item.Qty
	}
	return total
}

// ItemIds returns a slice of Ids or Items in the order
func (o Order) ItemIds() []string {

	var orderItemIds []string

	for i := 0; i <= len(o.Items)-1; i++ {

		item := o.Items[i]

		orderItemIds = append(orderItemIds, item.ItemID)

	}

	return orderItemIds
}

/* //////////////////////////
// CUSTOM ERRORS
*/ //////////////////////////

// ErrProcessOrder represents a process order error
type ErrProcessOrder struct {
	message string
}

// NewErrProcessOrder constructor
func NewErrProcessOrder(message string) *ErrProcessOrder {
	return &ErrProcessOrder{
		message: message,
	}
}

func (e *ErrProcessOrder) Error() string {
	return e.message
}

// ErrUpdateOrderStatus represents a process order error
type ErrUpdateOrderStatus struct {
	message string
}

// NewErrUpdateOrderStatus constructor
func NewErrUpdateOrderStatus(message string) *ErrUpdateOrderStatus {
	return &ErrUpdateOrderStatus{
		message: message,
	}
}

func (e *ErrUpdateOrderStatus) Error() string {
	return e.message
}
