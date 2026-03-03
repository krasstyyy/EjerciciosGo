// Package account provides a deliberately SHALLOW bank account module.
//
// This is the starting point for the coding dojo. Every field is exposed
// through getters and setters, and all business logic lives in free functions
// outside the struct. The caller must interrogate the account, make every
// decision, and mutate its state directly.
//
// Your job: refactor this through three rounds into a DEEP module.
package account

import (
	"errors"
	"fmt"
)

// Account is a shallow data holder — all fields are accessed through
// getters and setters, and all logic lives outside the struct.
type Account struct {
	owner          string
	balance        float64
	overdraftLimit float64
	frozen         bool
	transactions   []Transaction
}

type Transaction struct {
	name        string
	Amount      float64
	Balance     float64
	Description string
	Type        string
}

// New creates a new account with the given owner and overdraft limit.
func New(owner string, overdraftLimit float64) *Account {
	return &Account{
		owner:          owner,
		overdraftLimit: overdraftLimit,
	}
}

func (a *Account) Deposit(amount float64) error {
	if !a.frozen {
		if amount <= 0 {
			return fmt.Errorf("invalid amount: %.2f", amount)
		} else {
			a.balance += amount
			a.transactions = append(a.transactions, Transaction{
				name:        a.owner,
				Amount:      amount,
				Balance:     a.balance,
				Description: "deposit",
				Type:        "deposit",
			})
		}
	} else {
		return errors.New("account is frozen")
	}

	return nil
}

func (a *Account) Balance() float64 {
	return a.balance
}

func (a *Account) Withdraw(amount float64) error {
	if a.frozen {
		return errors.New("account is frozen")
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %.2f", amount)
	}
	if a.balance-amount < -a.overdraftLimit {
		return errors.New("insufficient funds")
	}
	a.balance -= amount
	a.transactions = append(a.transactions, Transaction{
		name:        a.owner,
		Amount:      amount,
		Balance:     a.balance,
		Description: "withdraw",
		Type:        "withdrawal",
	})
	return nil
}

func (a *Account) Transfer(to *Account, amount float64) error {
	if a.frozen || to.frozen {
		return errors.New("account is frozen")
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %.2f", amount)
	}
	if a.balance-amount < -a.overdraftLimit {
		return errors.New("insufficient funds")
	}
	to.balance += amount
	a.balance -= amount
	a.transactions = append(a.transactions, Transaction{
		name:        a.owner,
		Amount:      amount,
		Balance:     a.balance,
		Description: fmt.Sprintf("transfer to %s", to.owner),
		Type:        "transfer_out",
	})
	to.transactions = append(to.transactions, Transaction{
		name:        to.owner,
		Amount:      amount,
		Balance:     to.balance,
		Description: fmt.Sprintf("transfer from %s", a.owner),
		Type:        "transfer_in",
	})
	return nil
}

func (a *Account) Freeze() error {
	if a.frozen {
		return errors.New("account is already frozen")
	}
	a.frozen = true
	return nil
}

func (a *Account) Unfreeze() error {
	if !a.frozen {
		return errors.New("account is not frozen")
	}
	a.frozen = false
	return nil
}

func (a *Account) Transactions() []Transaction {

	return a.transactions
}
