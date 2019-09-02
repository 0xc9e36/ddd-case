package domain

import "errors"

type CustomerRepository interface {
	Store(customer Customer) error
	FindByID(ID int) Customer
}

type ItemRepository interface {
	Store(item Item) error
	FindByID(ID int) Item
}

type OrderRepository interface {
	Store(order Order) error
	FindByID(ID int) Order
}

type Customer struct {
	ID   int
	Name string
}

type Item struct {
	ID        int
	Name      string
	Value     float64
	Available bool
}

type Order struct {
	ID       int
	Customer Customer
	Items    []Item
}

func (o *Order) value() float64 {
	sum := 0.0
	for i := range o.Items {
		sum += o.Items[i].Value
	}
	return sum
}

//订单中添加订单项
func (o *Order) Add(item Item) error {

	if !item.Available {
		return errors.New("订单项不可添加!")
	}

	if o.value()+item.Value > 250.00 {
		return errors.New("订单总金额不能超过 250.00!")
	}

	o.Items = append(o.Items, item)
	return nil
}
