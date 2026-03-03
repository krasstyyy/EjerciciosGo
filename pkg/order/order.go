// Package order provides a deliberately SHALLOW order processing module.
//
// This is the bonus exercise for fast finishers. All fields are exposed
// through getters and setters, and all business logic lives in free functions
// that reach into every struct. Information leaks across module boundaries.
//
// Your job: refactor this through three rounds into DEEP modules.
package order

import (
	"errors"
	"fmt"
)

// --- Shallow data holders ---

// Item represents a line item with a SKU and quantity.
type Item struct {
	sku      string
	quantity int
}

type Event struct {
	Type string
}

func NewItem(sku string, quantity int) *Item {
	return &Item{sku: sku, quantity: quantity}
}

// Cart holds items before checkout.
type Cart struct {
	customerID string
	items      []*Item
}

func NewCart(customerID string) *Cart {
	return &Cart{customerID: customerID}
}

func (c *Cart) AddItem(item *Item) {
	c.items = append(c.items, item)
}

// Inventory tracks stock levels per SKU.
type Inventory struct {
	stock map[string]int
}

func NewInventory() *Inventory {
	return &Inventory{stock: make(map[string]int)}
}

func (inv *Inventory) GetStock(sku string) int {
	return inv.stock[sku]
}

func (inv *Inventory) SetStock(sku string, qty int) {
	inv.stock[sku] = qty
}

// Pricer holds unit prices per SKU.
type Pricer struct {
	prices map[string]float64
}

func NewPricer() *Pricer {
	return &Pricer{prices: make(map[string]float64)}
}

func (p *Pricer) GetPrice(sku string) float64 {
	return p.prices[sku]
}

func (p *Pricer) SetPrice(sku string, price float64) {
	p.prices[sku] = price
}

// Order tracks the lifecycle of a purchase.
type Order struct {
	id         string
	customerID string
	items      []*Item
	total      float64
	status     string // "pending", "paid", "shipped", "cancelled"
	address    string
	events     []Event
}

// --- Free functions: all business logic lives outside the structs ---

// PlaceOrder validates stock, calculates the total, reserves inventory, and
// creates a pending order. The caller must pass in every dependency and the
// function reaches into every struct's internals.
func PlaceOrder(cart *Cart, inventory *Inventory, pricer *Pricer) (*Order, error) {
	if len(cart.items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate stock for every item
	for _, item := range cart.items {
		if inventory.GetStock(item.sku) < item.quantity {
			return nil, fmt.Errorf("insufficient stock for %s", item.sku)
		}
	}

	// Calculate total
	total := 0.0
	for _, item := range cart.items {
		price := pricer.GetPrice(item.sku)
		if price <= 0 {
			return nil, fmt.Errorf("no price for %s", item.sku)
		}
		total += price * float64(item.quantity)
	}

	// Reserve inventory
	for _, item := range cart.items {
		inventory.SetStock(item.sku, inventory.GetStock(item.sku)-item.quantity)
	}

	return &Order{
		id:         fmt.Sprintf("ORD-%s-%d", cart.customerID, len(cart.items)),
		customerID: cart.customerID,
		items:      cart.items,
		total:      total,
		status:     "pending",
		events:     []Event{{Type: "created"}},
	}, nil
}

// Pay marks an order as paid. The caller must check status manually.
func (order *Order) Pay(amount float64) error {
	if order.status != "pending" {
		return fmt.Errorf("cannot pay order in status %q", order.status)
	}
	if amount < order.total {
		return fmt.Errorf("payment %.2f is less than total %.2f", amount, order.total)
	}
	order.status = "paid"
	order.events = append(order.events, Event{Type: "paid"})
	return nil
}

// Ship marks a paid order as shipped. The caller must check status and provide
// an address.
func (order *Order) Ship(address string) error {
	if order.status != "paid" {
		return fmt.Errorf("cannot ship order in status %q", order.status)
	}
	if address == "" {
		return errors.New("address is required")
	}
	order.address = address
	order.status = "shipped"
	order.events = append(order.events, Event{Type: "shipped"})
	return nil
}

// Cancel cancels an order and restores inventory. The caller must know which
// statuses are cancellable and manually restore stock.
func (order *Order) Cancel(inventory *Inventory) error {
	if order.status == "shipped" {
		return errors.New("cannot cancel a shipped order")
	}
	if order.status == "cancelled" {
		return errors.New("order is already cancelled")
	}

	// Restore inventory
	for _, item := range order.items {
		inventory.SetStock(item.sku, inventory.GetStock(item.sku)+item.quantity)
	}

	order.status = "cancelled"
	order.events = append(order.events, Event{Type: "cancelled"})
	return nil
}

func (cart *Cart) Checkout(inventory *Inventory, pricer *Pricer) (*Order, error) {

	return PlaceOrder(cart, inventory, pricer)
}

func (order *Order) Total() float64 {

	return order.total
}

func (order *Order) Status() string {

	return order.status
}

func (inventory *Inventory) Stock(sku string) int {

	return inventory.stock[sku]
}

func (order *Order) Events() []Event {
	return order.events
}
