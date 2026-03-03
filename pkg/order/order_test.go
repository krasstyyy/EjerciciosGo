package order

import (
	"testing"
)

// =============================================================================
// Helpers
// =============================================================================

func setupInventoryAndPricer() (*Inventory, *Pricer) {
	inv := NewInventory()
	inv.SetStock("TSHIRT", 10)
	inv.SetStock("MUG", 5)

	pricer := NewPricer()
	pricer.SetPrice("TSHIRT", 25.00)
	pricer.SetPrice("MUG", 12.50)

	return inv, pricer
}

// =============================================================================
// Round 0 — Tests for the shallow "before" code
// =============================================================================

/*func TestPlaceOrder(t *testing.T) {
	tests := []struct {
		name      string
		items     []*Item
		stock     map[string]int
		prices    map[string]float64
		wantErr   bool
		wantTotal float64
	}{
		{
			name:      "valid order",
			items:     []*Item{NewItem("TSHIRT", 2), NewItem("MUG", 1)},
			stock:     map[string]int{"TSHIRT": 10, "MUG": 5},
			prices:    map[string]float64{"TSHIRT": 25.00, "MUG": 12.50},
			wantTotal: 62.50,
		},
		{
			name:    "empty cart",
			items:   nil,
			stock:   map[string]int{},
			prices:  map[string]float64{},
			wantErr: true,
		},
		{
			name:    "insufficient stock",
			items:   []*Item{NewItem("TSHIRT", 20)},
			stock:   map[string]int{"TSHIRT": 10},
			prices:  map[string]float64{"TSHIRT": 25.00},
			wantErr: true,
		},
		{
			name:    "missing price",
			items:   []*Item{NewItem("TSHIRT", 1)},
			stock:   map[string]int{"TSHIRT": 10},
			prices:  map[string]float64{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart("customer-1")
			for _, item := range tt.items {
				cart.AddItem(item)
			}

			inv := NewInventory()
			for sku, qty := range tt.stock {
				inv.SetStock(sku, qty)
			}

			pricer := NewPricer()
			for sku, price := range tt.prices {
				pricer.SetPrice(sku, price)
			}

			ord, err := PlaceOrder(cart, inv, pricer)

			if (err != nil) != tt.wantErr {
				t.Fatalf("PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if ord.GetTotal() != tt.wantTotal {
					t.Errorf("expected total %.2f, got %.2f", tt.wantTotal, ord.GetTotal())
				}
				if ord.GetStatus() != "pending" {
					t.Errorf("expected status pending, got %s", ord.GetStatus())
				}
			}
		})
	}
}

func TestPlaceOrderReservesInventory(t *testing.T) {
	tests := []struct {
		name      string
		sku       string
		quantity  int
		initial   int
		wantStock int
	}{
		{
			name:      "reserves stock on placement",
			sku:       "TSHIRT",
			quantity:  3,
			initial:   10,
			wantStock: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			inv.SetStock(tt.sku, tt.initial)

			cart := NewCart("customer-1")
			cart.AddItem(NewItem(tt.sku, tt.quantity))

			_, err := PlaceOrder(cart, inv, pricer)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if inv.GetStock(tt.sku) != tt.wantStock {
				t.Errorf("expected stock %d, got %d", tt.wantStock, inv.GetStock(tt.sku))
			}
		})
	}
}

func TestPay(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		amount  float64
		total   float64
		wantErr bool
	}{
		{
			name:   "valid payment",
			status: "pending",
			amount: 50.00,
			total:  50.00,
		},
		{
			name:   "overpayment accepted",
			status: "pending",
			amount: 100.00,
			total:  50.00,
		},
		{
			name:    "underpayment rejected",
			status:  "pending",
			amount:  10.00,
			total:   50.00,
			wantErr: true,
		},
		{
			name:    "already paid",
			status:  "paid",
			amount:  50.00,
			total:   50.00,
			wantErr: true,
		},
		{
			name:    "already shipped",
			status:  "shipped",
			amount:  50.00,
			total:   50.00,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ord := &Order{status: tt.status, total: tt.total}

			err := Pay(ord, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Pay() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && ord.GetStatus() != "paid" {
				t.Errorf("expected status paid, got %s", ord.GetStatus())
			}
		})
	}
}

func TestShip(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		address     string
		wantErr     bool
		wantAddress string
	}{
		{
			name:        "valid shipment",
			status:      "paid",
			address:     "123 Main St",
			wantAddress: "123 Main St",
		},
		{
			name:    "not yet paid",
			status:  "pending",
			address: "123 Main St",
			wantErr: true,
		},
		{
			name:    "empty address",
			status:  "paid",
			address: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ord := &Order{status: tt.status}

			err := Ship(ord, tt.address)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Ship() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if ord.GetStatus() != "shipped" {
					t.Errorf("expected status shipped, got %s", ord.GetStatus())
				}
				if ord.GetAddress() != tt.wantAddress {
					t.Errorf("expected address %q, got %q", tt.wantAddress, ord.GetAddress())
				}
			}
		})
	}
}

func TestCancel(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		wantErr   bool
		wantStock int
	}{
		{
			name:      "cancel pending order",
			status:    "pending",
			wantStock: 10,
		},
		{
			name:      "cancel paid order",
			status:    "paid",
			wantStock: 10,
		},
		{
			name:    "cannot cancel shipped order",
			status:  "shipped",
			wantErr: true,
		},
		{
			name:    "cannot cancel already cancelled",
			status:  "cancelled",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := NewInventory()
			inv.SetStock("TSHIRT", 7)

			ord := &Order{
				status: tt.status,
				items:  []*Item{NewItem("TSHIRT", 3)},
			}

			err := Cancel(ord, inv)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if ord.GetStatus() != "cancelled" {
					t.Errorf("expected status cancelled, got %s", ord.GetStatus())
				}
				if inv.GetStock("TSHIRT") != tt.wantStock {
					t.Errorf("expected stock %d, got %d", tt.wantStock, inv.GetStock("TSHIRT"))
				}
			}
		})
	}
}*/

// =============================================================================
// Round 1 — Tell, Don't Ask
//
// Move logic into methods. The cart handles its own checkout, the order manages
// its own state transitions.
//
// cart.Checkout(inventory, pricer) → *Order
// order.Pay(amount) → error
// order.Ship(address) → error
// order.Cancel(inventory) → error
// =============================================================================

func TestRound1_Checkout(t *testing.T) {
	tests := []struct {
		name      string
		items     []*Item
		wantErr   bool
		wantTotal float64
		wantStock map[string]int
	}{
		{
			name:      "valid checkout",
			items:     []*Item{NewItem("TSHIRT", 2), NewItem("MUG", 1)},
			wantTotal: 62.50,
			wantStock: map[string]int{"TSHIRT": 8, "MUG": 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			cart := NewCart("customer-1")
			for _, item := range tt.items {
				cart.AddItem(item)
			}

			ord, err := cart.Checkout(inv, pricer)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if ord.Total() != tt.wantTotal {
					t.Errorf("expected total %.2f, got %.2f", tt.wantTotal, ord.Total())
				}
				if ord.Status() != "pending" {
					t.Errorf("expected status pending, got %s", ord.Status())
				}
				for sku, want := range tt.wantStock {
					if inv.Stock(sku) != want {
						t.Errorf("expected %s stock %d, got %d", sku, want, inv.Stock(sku))
					}
				}
			}
		})
	}
}

func TestRound1_Lifecycle(t *testing.T) {
	tests := []struct {
		name       string
		pay        float64
		ship       string
		wantStatus string
	}{
		{
			name:       "full lifecycle: pending → paid → shipped",
			pay:        25.00,
			ship:       "123 Main St",
			wantStatus: "shipped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			cart := NewCart("customer-1")
			cart.AddItem(NewItem("TSHIRT", 1))

			ord, _ := cart.Checkout(inv, pricer)

			if err := ord.Pay(tt.pay); err != nil {
				t.Fatalf("Pay() unexpected error: %v", err)
			}
			if err := ord.Ship(tt.ship); err != nil {
				t.Fatalf("Ship() unexpected error: %v", err)
			}
			if ord.Status() != tt.wantStatus {
				t.Errorf("expected status %s, got %s", tt.wantStatus, ord.Status())
			}
		})
	}
}

func TestRound1_CancelRestoresStock(t *testing.T) {
	tests := []struct {
		name            string
		sku             string
		quantity        int
		wantStockAfter  int
		wantStockCancel int
	}{
		{
			name:            "cancel restores inventory",
			sku:             "MUG",
			quantity:        3,
			wantStockAfter:  2,
			wantStockCancel: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			cart := NewCart("customer-1")
			cart.AddItem(NewItem(tt.sku, tt.quantity))

			ord, _ := cart.Checkout(inv, pricer)
			if inv.Stock(tt.sku) != tt.wantStockAfter {
				t.Fatalf("expected stock %d after checkout, got %d", tt.wantStockAfter, inv.Stock(tt.sku))
			}

			if err := ord.Cancel(inv); err != nil {
				t.Fatalf("Cancel() unexpected error: %v", err)
			}
			if inv.Stock(tt.sku) != tt.wantStockCancel {
				t.Errorf("expected stock %d after cancel, got %d", tt.wantStockCancel, inv.Stock(tt.sku))
			}
		})
	}
}

// =============================================================================
// Round 2 — Deep Modules: Event Log
//
// Add an event log to Order. Each state transition records an Event with
// timestamp, type, and description.
//
// type Event struct {
//     Timestamp   time.Time
//     Type        string   // "created", "paid", "shipped", "cancelled"
//     Description string
// }
// =============================================================================

func TestRound2_CheckoutEvent(t *testing.T) {
	tests := []struct {
		name       string
		items      []*Item
		wantEvents int
		wantType   string
	}{
		{
			name:       "checkout records created event",
			items:      []*Item{NewItem("TSHIRT", 1)},
			wantEvents: 1,
			wantType:   "created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			cart := NewCart("customer-1")
			for _, item := range tt.items {
				cart.AddItem(item)
			}

			ord, _ := cart.Checkout(inv, pricer)

			events := ord.Events()
			if len(events) != tt.wantEvents {
				t.Fatalf("expected %d event(s), got %d", tt.wantEvents, len(events))
			}
			if events[0].Type != tt.wantType {
				t.Errorf("expected type %s, got %s", tt.wantType, events[0].Type)
			}
		})
	}
}

func TestRound2_LifecycleEvents(t *testing.T) {
	tests := []struct {
		name      string
		pay       float64
		ship      string
		wantTypes []string
	}{
		{
			name:      "full lifecycle records all events",
			pay:       25.00,
			ship:      "123 Main St",
			wantTypes: []string{"created", "paid", "shipped"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			cart := NewCart("customer-1")
			cart.AddItem(NewItem("TSHIRT", 1))

			ord, _ := cart.Checkout(inv, pricer)
			_ = ord.Pay(tt.pay)
			_ = ord.Ship(tt.ship)

			events := ord.Events()
			if len(events) != len(tt.wantTypes) {
				t.Fatalf("expected %d events, got %d", len(tt.wantTypes), len(events))
			}
			for i, want := range tt.wantTypes {
				if events[i].Type != want {
					t.Errorf("event[%d]: expected type %s, got %s", i, want, events[i].Type)
				}
			}
		})
	}
}

func TestRound2_FailedTransitionNoEvent(t *testing.T) {
	tests := []struct {
		name       string
		action     string // "ship" on unpaid order
		wantEvents int
	}{
		{
			name:       "failed ship does not record event",
			action:     "ship",
			wantEvents: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, pricer := setupInventoryAndPricer()
			cart := NewCart("customer-1")
			cart.AddItem(NewItem("TSHIRT", 1))

			ord, _ := cart.Checkout(inv, pricer)
			_ = ord.Ship("123 Main St") // should fail — not paid yet

			events := ord.Events()
			if len(events) != tt.wantEvents {
				t.Errorf("expected %d event(s), got %d", tt.wantEvents, len(events))
			}
		})
	}
}
