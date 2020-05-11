package main

import (
	"math"

	"golang.org/x/text/currency"
)

// Money represents a real world money
type Money struct {
	currency.Amount
	value int
}

// NewMoney creates a new instance of the money class
func NewMoney(currency currency.Unit, amount int) Money {
	return Money{currency.Amount(amount), amount}
}

// NewEUR creates a new EURO money
func NewEUR(amount int) Money {
	return NewMoney(currency.EUR, amount)
}

// Multiply multiplies the money amount by a float and rounds the value up
func (money Money) Multiply(value float64) Money {
	newAmount := int(math.Round(float64(money.value) * value))
	return money.AddAmount(newAmount)
}

// AddAmount increments the current money by an amount
func (money Money) AddAmount(amount int) (result Money) {
	return NewMoney(money.Currency(), money.value+amount)
}

// Value returns the value of the money in the base units
func (money Money) Value() int {
	return money.value
}
