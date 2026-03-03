package account

import (
	"testing"
)

// =============================================================================
// Round 0 — Tests for the shallow "before" code
// =============================================================================

func TestNew(t *testing.T) {
	acc := New("Alice", 100.0)

	if acc.GetOwner() != "Alice" {
		t.Errorf("expected owner Alice, got %s", acc.GetOwner())
	}
	if acc.GetBalance() != 0 {
		t.Errorf("expected balance 0, got %.2f", acc.GetBalance())
	}
	if acc.GetOverdraftLimit() != 100.0 {
		t.Errorf("expected overdraft limit 100, got %.2f", acc.GetOverdraftLimit())
	}
	if acc.GetIsFrozen() {
		t.Error("expected account to not be frozen")
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name        string
		initial     float64
		frozen      bool
		amount      float64
		wantErr     bool
		wantBalance float64
	}{
		{
			name:        "valid deposit",
			amount:      200.0,
			wantBalance: 200.0,
		},
		{
			name:    "zero amount",
			amount:  0,
			wantErr: true,
		},
		{
			name:    "negative amount",
			amount:  -50,
			wantErr: true,
		},
		{
			name:    "frozen account",
			frozen:  true,
			amount:  100.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New("Alice", 0)
			if tt.initial > 0 {
				_ = Deposit(acc, tt.initial)
			}
			if tt.frozen {
				_ = Freeze(acc)
			}

			err := Deposit(acc, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && acc.GetBalance() != tt.wantBalance {
				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, acc.GetBalance())
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name           string
		initial        float64
		overdraftLimit float64
		frozen         bool
		amount         float64
		wantErr        bool
		wantBalance    float64
	}{
		{
			name:        "valid withdrawal",
			initial:     500.0,
			amount:      200.0,
			wantBalance: 300.0,
		},
		{
			name:           "within overdraft limit",
			initial:        50.0,
			overdraftLimit: 100.0,
			amount:         120.0,
			wantBalance:    -70.0,
		},
		{
			name:           "exceeds overdraft limit",
			initial:        50.0,
			overdraftLimit: 100.0,
			amount:         200.0,
			wantErr:        true,
			wantBalance:    50.0,
		},
		{
			name:    "zero amount",
			amount:  0,
			wantErr: true,
		},
		{
			name:    "negative amount",
			amount:  -10,
			wantErr: true,
		},
		{
			name:        "frozen account",
			initial:     100.0,
			frozen:      true,
			amount:      50.0,
			wantErr:     true,
			wantBalance: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New("Alice", tt.overdraftLimit)
			if tt.initial > 0 {
				_ = Deposit(acc, tt.initial)
			}
			if tt.frozen {
				_ = Freeze(acc)
			}

			err := Withdraw(acc, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && acc.GetBalance() != tt.wantBalance {
				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, acc.GetBalance())
			}
			if tt.wantErr && tt.initial > 0 && acc.GetBalance() != tt.wantBalance {
				t.Errorf("balance should be unchanged, got %.2f", acc.GetBalance())
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	tests := []struct {
		name            string
		fromBalance     float64
		toBalance       float64
		fromFrozen      bool
		toFrozen        bool
		amount          float64
		wantErr         bool
		wantFromBalance float64
		wantToBalance   float64
	}{
		{
			name:            "valid transfer",
			fromBalance:     500.0,
			amount:          200.0,
			wantFromBalance: 300.0,
			wantToBalance:   200.0,
		},
		{
			name:            "insufficient funds",
			fromBalance:     100.0,
			amount:          200.0,
			wantErr:         true,
			wantFromBalance: 100.0,
			wantToBalance:   0,
		},
		{
			name:    "zero amount",
			amount:  0,
			wantErr: true,
		},
		{
			name:    "negative amount",
			amount:  -10,
			wantErr: true,
		},
		{
			name:        "source account frozen",
			fromBalance: 500.0,
			fromFrozen:  true,
			amount:      100.0,
			wantErr:     true,
		},
		{
			name:        "destination account frozen",
			fromBalance: 500.0,
			toFrozen:    true,
			amount:      100.0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from := New("Alice", 0)
			to := New("Bob", 0)
			if tt.fromBalance > 0 {
				_ = Deposit(from, tt.fromBalance)
			}
			if tt.toBalance > 0 {
				_ = Deposit(to, tt.toBalance)
			}
			if tt.fromFrozen {
				_ = Freeze(from)
			}
			if tt.toFrozen {
				_ = Freeze(to)
			}

			err := Transfer(from, to, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if from.GetBalance() != tt.wantFromBalance {
					t.Errorf("expected sender balance %.2f, got %.2f", tt.wantFromBalance, from.GetBalance())
				}
				if to.GetBalance() != tt.wantToBalance {
					t.Errorf("expected receiver balance %.2f, got %.2f", tt.wantToBalance, to.GetBalance())
				}
			}
		})
	}
}

func TestFreeze(t *testing.T) {
	tests := []struct {
		name       string
		frozen     bool
		wantErr    bool
		wantFrozen bool
	}{
		{
			name:       "freeze unfrozen account",
			wantFrozen: true,
		},
		{
			name:       "freeze already frozen account",
			frozen:     true,
			wantErr:    true,
			wantFrozen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New("Alice", 0)
			if tt.frozen {
				_ = Freeze(acc)
			}

			err := Freeze(acc)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Freeze() error = %v, wantErr %v", err, tt.wantErr)
			}
			if acc.GetIsFrozen() != tt.wantFrozen {
				t.Errorf("expected frozen=%v, got %v", tt.wantFrozen, acc.GetIsFrozen())
			}
		})
	}
}

func TestUnfreeze(t *testing.T) {
	tests := []struct {
		name       string
		frozen     bool
		wantErr    bool
		wantFrozen bool
	}{
		{
			name:   "unfreeze frozen account",
			frozen: true,
		},
		{
			name:    "unfreeze non-frozen account",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New("Alice", 0)
			if tt.frozen {
				_ = Freeze(acc)
			}

			err := Unfreeze(acc)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Unfreeze() error = %v, wantErr %v", err, tt.wantErr)
			}
			if acc.GetIsFrozen() != tt.wantFrozen {
				t.Errorf("expected frozen=%v, got %v", tt.wantFrozen, acc.GetIsFrozen())
			}
		})
	}
}

// =============================================================================
// Round 1 — Tell, Don't Ask
//
// After refactoring, uncomment these tests. The operations become methods
// on Account: acc.Withdraw(100), from.Transfer(to, 50), etc.
// =============================================================================

// func TestRound1_Deposit(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		amount      float64
// 		wantErr     bool
// 		wantBalance float64
// 	}{
// 		{
// 			name:        "valid deposit",
// 			amount:      300.0,
// 			wantBalance: 300.0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			acc := New("Alice", 0)
//
// 			err := acc.Deposit(tt.amount)
//
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if !tt.wantErr && acc.Balance() != tt.wantBalance {
// 				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, acc.Balance())
// 			}
// 		})
// 	}
// }

// func TestRound1_Withdraw(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		initial        float64
// 		overdraftLimit float64
// 		amount         float64
// 		wantErr        bool
// 		wantBalance    float64
// 	}{
// 		{
// 			name:        "valid withdrawal",
// 			initial:     200.0,
// 			amount:      50.0,
// 			wantBalance: 150.0,
// 		},
// 		{
// 			name:           "within overdraft limit",
// 			initial:        50.0,
// 			overdraftLimit: 100.0,
// 			amount:         120.0,
// 			wantBalance:    -70.0,
// 		},
// 		{
// 			name:           "exceeds overdraft limit",
// 			initial:        50.0,
// 			overdraftLimit: 100.0,
// 			amount:         200.0,
// 			wantErr:        true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			acc := New("Alice", tt.overdraftLimit)
// 			_ = acc.Deposit(tt.initial)
//
// 			err := acc.Withdraw(tt.amount)
//
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if !tt.wantErr && acc.Balance() != tt.wantBalance {
// 				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, acc.Balance())
// 			}
// 		})
// 	}
// }

// func TestRound1_Transfer(t *testing.T) {
// 	tests := []struct {
// 		name            string
// 		fromBalance     float64
// 		amount          float64
// 		wantErr         bool
// 		wantFromBalance float64
// 		wantToBalance   float64
// 	}{
// 		{
// 			name:            "valid transfer",
// 			fromBalance:     500.0,
// 			amount:          200.0,
// 			wantFromBalance: 300.0,
// 			wantToBalance:   200.0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			from := New("Alice", 0)
// 			to := New("Bob", 0)
// 			_ = from.Deposit(tt.fromBalance)
//
// 			err := from.Transfer(to, tt.amount)
//
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if !tt.wantErr {
// 				if from.Balance() != tt.wantFromBalance {
// 					t.Errorf("expected sender balance %.2f, got %.2f", tt.wantFromBalance, from.Balance())
// 				}
// 				if to.Balance() != tt.wantToBalance {
// 					t.Errorf("expected receiver balance %.2f, got %.2f", tt.wantToBalance, to.Balance())
// 				}
// 			}
// 		})
// 	}
// }

// func TestRound1_Freeze(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		freeze      bool
// 		unfreeze    bool
// 		deposit     float64
// 		wantErr     bool
// 		wantBalance float64
// 	}{
// 		{
// 			name:    "frozen account blocks deposit",
// 			freeze:  true,
// 			deposit: 10.0,
// 			wantErr: true,
// 		},
// 		{
// 			name:        "unfreeze re-enables deposit",
// 			freeze:      true,
// 			unfreeze:    true,
// 			deposit:     10.0,
// 			wantBalance: 110.0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			acc := New("Alice", 0)
// 			_ = acc.Deposit(100.0)
//
// 			if tt.freeze {
// 				acc.Freeze()
// 			}
// 			if tt.unfreeze {
// 				acc.Unfreeze()
// 			}
//
// 			err := acc.Deposit(tt.deposit)
//
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if !tt.wantErr && acc.Balance() != tt.wantBalance {
// 				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, acc.Balance())
// 			}
// 		})
// 	}
// }

// =============================================================================
// Round 2 — Deep Module: Transaction History
//
// After adding transaction logging, uncomment these tests.
// =============================================================================

// func TestRound2_DepositTransaction(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		amount      float64
// 		wantType    string
// 		wantAmount  float64
// 		wantBalance float64
// 	}{
// 		{
// 			name:        "deposit records transaction",
// 			amount:      500.0,
// 			wantType:    "deposit",
// 			wantAmount:  500.0,
// 			wantBalance: 500.0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			acc := New("Alice", 0)
// 			_ = acc.Deposit(tt.amount)
//
// 			txns := acc.Transactions()
// 			if len(txns) != 1 {
// 				t.Fatalf("expected 1 transaction, got %d", len(txns))
// 			}
// 			if txns[0].Type != tt.wantType {
// 				t.Errorf("expected type %s, got %s", tt.wantType, txns[0].Type)
// 			}
// 			if txns[0].Amount != tt.wantAmount {
// 				t.Errorf("expected amount %.2f, got %.2f", tt.wantAmount, txns[0].Amount)
// 			}
// 			if txns[0].Balance != tt.wantBalance {
// 				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, txns[0].Balance)
// 			}
// 		})
// 	}
// }

// func TestRound2_WithdrawTransaction(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		deposit     float64
// 		withdraw    float64
// 		wantType    string
// 		wantAmount  float64
// 		wantBalance float64
// 	}{
// 		{
// 			name:        "withdraw records transaction",
// 			deposit:     500.0,
// 			withdraw:    200.0,
// 			wantType:    "withdrawal",
// 			wantAmount:  200.0,
// 			wantBalance: 300.0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			acc := New("Alice", 0)
// 			_ = acc.Deposit(tt.deposit)
// 			_ = acc.Withdraw(tt.withdraw)
//
// 			txns := acc.Transactions()
// 			if len(txns) != 2 {
// 				t.Fatalf("expected 2 transactions, got %d", len(txns))
// 			}
// 			tx := txns[1]
// 			if tx.Type != tt.wantType {
// 				t.Errorf("expected type %s, got %s", tt.wantType, tx.Type)
// 			}
// 			if tx.Amount != tt.wantAmount {
// 				t.Errorf("expected amount %.2f, got %.2f", tt.wantAmount, tx.Amount)
// 			}
// 			if tx.Balance != tt.wantBalance {
// 				t.Errorf("expected balance %.2f, got %.2f", tt.wantBalance, tx.Balance)
// 			}
// 		})
// 	}
// }

// func TestRound2_TransferTransactions(t *testing.T) {
// 	tests := []struct {
// 		name            string
// 		deposit         float64
// 		transfer        float64
// 		wantFromType    string
// 		wantFromBalance float64
// 		wantFromDesc    string
// 		wantToType      string
// 		wantToBalance   float64
// 		wantToDesc      string
// 	}{
// 		{
// 			name:            "transfer records both sides",
// 			deposit:         500.0,
// 			transfer:        150.0,
// 			wantFromType:    "transfer_out",
// 			wantFromBalance: 350.0,
// 			wantFromDesc:    "transfer to Bob",
// 			wantToType:      "transfer_in",
// 			wantToBalance:   150.0,
// 			wantToDesc:      "transfer from Alice",
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			from := New("Alice", 0)
// 			to := New("Bob", 0)
// 			_ = from.Deposit(tt.deposit)
// 			_ = from.Transfer(to, tt.transfer)
//
// 			fromTxns := from.Transactions()
// 			if len(fromTxns) != 2 {
// 				t.Fatalf("expected 2 sender transactions, got %d", len(fromTxns))
// 			}
// 			if fromTxns[1].Type != tt.wantFromType {
// 				t.Errorf("sender: expected type %s, got %s", tt.wantFromType, fromTxns[1].Type)
// 			}
// 			if fromTxns[1].Amount != tt.transfer {
// 				t.Errorf("sender: expected amount %.2f, got %.2f", tt.transfer, fromTxns[1].Amount)
// 			}
// 			if fromTxns[1].Balance != tt.wantFromBalance {
// 				t.Errorf("sender: expected balance %.2f, got %.2f", tt.wantFromBalance, fromTxns[1].Balance)
// 			}
// 			if fromTxns[1].Description != tt.wantFromDesc {
// 				t.Errorf("sender: expected description %q, got %q", tt.wantFromDesc, fromTxns[1].Description)
// 			}
//
// 			toTxns := to.Transactions()
// 			if len(toTxns) != 1 {
// 				t.Fatalf("expected 1 receiver transaction, got %d", len(toTxns))
// 			}
// 			if toTxns[0].Type != tt.wantToType {
// 				t.Errorf("receiver: expected type %s, got %s", tt.wantToType, toTxns[0].Type)
// 			}
// 			if toTxns[0].Amount != tt.transfer {
// 				t.Errorf("receiver: expected amount %.2f, got %.2f", tt.transfer, toTxns[0].Amount)
// 			}
// 			if toTxns[0].Balance != tt.wantToBalance {
// 				t.Errorf("receiver: expected balance %.2f, got %.2f", tt.wantToBalance, toTxns[0].Balance)
// 			}
// 			if toTxns[0].Description != tt.wantToDesc {
// 				t.Errorf("receiver: expected description %q, got %q", tt.wantToDesc, toTxns[0].Description)
// 			}
// 		})
// 	}
// }

// func TestRound2_FailedOperationNoTransaction(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		deposit  float64
// 		withdraw float64
// 		wantTxns int
// 	}{
// 		{
// 			name:     "failed withdraw does not record",
// 			deposit:  100.0,
// 			withdraw: 500.0,
// 			wantTxns: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			acc := New("Alice", 0)
// 			_ = acc.Deposit(tt.deposit)
// 			_ = acc.Withdraw(tt.withdraw)
//
// 			txns := acc.Transactions()
// 			if len(txns) != tt.wantTxns {
// 				t.Errorf("expected %d transaction(s), got %d", tt.wantTxns, len(txns))
// 			}
// 		})
// 	}
// }

