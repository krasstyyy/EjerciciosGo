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
}

// New creates a new account with the given owner and overdraft limit.
func New(owner string, overdraftLimit float64) *Account {
	return &Account{
		owner:          owner,
		overdraftLimit: overdraftLimit,
	}
}

func (a *Account) GetOwner() string {
	return a.owner
}

func (a *Account) GetBalance() float64 {
	return a.balance
}

func (a *Account) GetOverdraftLimit() float64 {
	return a.overdraftLimit
}

func (a *Account) GetIsFrozen() bool {
	return a.frozen
}

func (a *Account) SetOwner(owner string) {
	a.owner = owner
}

func (a *Account) SetBalance(balance float64) {
	a.balance = balance
}

func (a *Account) SetOverdraftLimit(limit float64) {
	a.overdraftLimit = limit
}

func (a *Account) SetIsFrozen(frozen bool) {
	a.frozen = frozen
}

// Withdraw removes money from the account.
// The caller must deal with: frozen checks, amount validation, overdraft logic.
func Withdraw(account *Account, amount float64) error {
	if account.GetIsFrozen() {
		return errors.New("account is frozen")
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %.2f", amount)
	}
	if account.GetBalance()-amount < -account.GetOverdraftLimit() {
		return errors.New("insufficient funds")
	}
	account.SetBalance(account.GetBalance() - amount)
	return nil
}

// Deposit adds money to the account.
// The caller must deal with: frozen checks, amount validation.
func Deposit(account *Account, amount float64) error {
	if account.GetIsFrozen() {
		return errors.New("account is frozen")
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %.2f", amount)
	}
	account.SetBalance(account.GetBalance() + amount)
	return nil
}

// Transfer moves money from one account to another.
// The caller must deal with: frozen checks on both accounts, amount validation,
// overdraft logic — reaching into the internals of two objects.
func Transfer(from, to *Account, amount float64) error {
	if from.GetIsFrozen() {
		return errors.New("source account is frozen")
	}
	if to.GetIsFrozen() {
		return errors.New("destination account is frozen")
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %.2f", amount)
	}
	if from.GetBalance()-amount < -from.GetOverdraftLimit() {
		return errors.New("insufficient funds")
	}
	from.SetBalance(from.GetBalance() - amount)
	to.SetBalance(to.GetBalance() + amount)
	return nil
}

// Freeze prevents any further operations on the account.
func Freeze(account *Account) error {
	if account.GetIsFrozen() {
		return errors.New("account is already frozen")
	}
	account.SetIsFrozen(true)
	return nil
}

// Unfreeze re-enables operations on the account.
func Unfreeze(account *Account) error {
	if !account.GetIsFrozen() {
		return errors.New("account is not frozen")
	}
	account.SetIsFrozen(false)
	return nil
}
