package main

import "fmt"

type PaymentMethod interface {
	pay(amount float64) string
}

type CreditCard struct {
	Number string
}

type UPI struct {
	ID string
}

func (c CreditCard) pay(amount float64) string {
	return fmt.Sprintf("Paid %f using Credit Card", amount)
}

func (upi UPI) pay(amount float64) string {
	return fmt.Sprintf("Paid %f using UPI", amount)
}

func processPayment(pay PaymentMethod, amount float64) string {
	return pay.pay(amount)
}

func main() {
	c := CreditCard{Number: "123456789"}
	upi := UPI{ID: "abc@pay"}

	fmt.Println(processPayment(c, 500.0))
	fmt.Println(processPayment(upi, 500.0))
}