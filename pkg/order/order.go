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

func NewItem(sku string, quantity int) *Item {
	return &Item{sku: sku, quantity: quantity}
}

func (i *Item) GetSKU() string {
	return i.sku
}

func (i *Item) GetQuantity() int {
	return i.quantity
}

func (i *Item) SetSKU(sku string) {
	i.sku = sku
}

func (i *Item) SetQuantity(qty int) {
	i.quantity = qty
}

// Cart holds items before checkout.
type Cart struct {
	customerID string
	items      []*Item
}

func NewCart(customerID string) *Cart {
	return &Cart{customerID: customerID}
}

func (c *Cart) GetCustomerID() string {
	return c.customerID
}

func (c *Cart) GetItems() []*Item {
	return c.items
}

func (c *Cart) SetItems(items []*Item) {
	c.items = items
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
}

func (o *Order) GetID() string {
	return o.id
}

func (o *Order) GetCustomerID() string {
	return o.customerID
}

func (o *Order) GetItems() []*Item {
	return o.items
}

func (o *Order) GetTotal() float64 {
	return o.total
}

func (o *Order) GetStatus() string {
	return o.status
}

func (o *Order) GetAddress() string {
	return o.address
}

func (o *Order) SetID(id string) {
	o.id = id
}

func (o *Order) SetStatus(s string) {
	o.status = s
}

func (o *Order) SetTotal(t float64) {
	o.total = t
}

func (o *Order) SetAddress(addr string) {
	o.address = addr
}

// --- Free functions: all business logic lives outside the structs ---

// PlaceOrder validates stock, calculates the total, reserves inventory, and
// creates a pending order. The caller must pass in every dependency and the
// function reaches into every struct's internals.
func PlaceOrder(cart *Cart, inventory *Inventory, pricer *Pricer) (*Order, error) {
	if len(cart.GetItems()) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate stock for every item
	for _, item := range cart.GetItems() {
		if inventory.GetStock(item.GetSKU()) < item.GetQuantity() {
			return nil, fmt.Errorf("insufficient stock for %s", item.GetSKU())
		}
	}

	// Calculate total
	total := 0.0
	for _, item := range cart.GetItems() {
		price := pricer.GetPrice(item.GetSKU())
		if price <= 0 {
			return nil, fmt.Errorf("no price for %s", item.GetSKU())
		}
		total += price * float64(item.GetQuantity())
	}

	// Reserve inventory
	for _, item := range cart.GetItems() {
		inventory.SetStock(item.GetSKU(), inventory.GetStock(item.GetSKU())-item.GetQuantity())
	}

	return &Order{
		id:         fmt.Sprintf("ORD-%s-%d", cart.GetCustomerID(), len(cart.GetItems())),
		customerID: cart.GetCustomerID(),
		items:      cart.GetItems(),
		total:      total,
		status:     "pending",
	}, nil
}

// Pay marks an order as paid. The caller must check status manually.
func Pay(order *Order, amount float64) error {
	if order.GetStatus() != "pending" {
		return fmt.Errorf("cannot pay order in status %q", order.GetStatus())
	}
	if amount < order.GetTotal() {
		return fmt.Errorf("payment %.2f is less than total %.2f", amount, order.GetTotal())
	}
	order.SetStatus("paid")
	return nil
}

// Ship marks a paid order as shipped. The caller must check status and provide
// an address.
func Ship(order *Order, address string) error {
	if order.GetStatus() != "paid" {
		return fmt.Errorf("cannot ship order in status %q", order.GetStatus())
	}
	if address == "" {
		return errors.New("address is required")
	}
	order.SetAddress(address)
	order.SetStatus("shipped")
	return nil
}

// Cancel cancels an order and restores inventory. The caller must know which
// statuses are cancellable and manually restore stock.
func Cancel(order *Order, inventory *Inventory) error {
	if order.GetStatus() == "shipped" {
		return errors.New("cannot cancel a shipped order")
	}
	if order.GetStatus() == "cancelled" {
		return errors.New("order is already cancelled")
	}

	// Restore inventory
	for _, item := range order.GetItems() {
		inventory.SetStock(item.GetSKU(), inventory.GetStock(item.GetSKU())+item.GetQuantity())
	}

	order.SetStatus("cancelled")
	return nil
}
