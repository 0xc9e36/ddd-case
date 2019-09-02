package usecases

import (
	"domain"
	"errors"
)

type User struct {
	ID       int
	IsAdmin  bool
	Customer domain.Customer
}

type UserRepository interface {
	Store(user User) error
	FindByID(ID int) User
}

type Item struct {
	ID    int
	Name  string
	Value float64
}

type Logger interface {
	Log(message string) error
}

type OrderInteractor struct {
	UserRepository  UserRepository
	OrderRepository domain.OrderRepository
	ItemRepository  domain.ItemRepository
	Logger          Logger
}

func (o *OrderInteractor) Items(userID, orderID int) ([]Item, error) {
	var items []Item
	user := o.UserRepository.FindByID(userID)
	order := o.OrderRepository.FindByID(orderID)

	if order.Customer.ID != user.Customer.ID {
		items = make([]Item, 0)
		err := errors.New("订单 ID 与用户 ID 不符合!")
		o.Logger.Log(err.Error())
		return items, err
	}

	items = make([]Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = Item{
			item.ID, item.Name, item.Value,
		}
	}
	return items, nil
}

func (o *OrderInteractor) Add(userID, orderID, itemID int) error {
	user := o.UserRepository.FindByID(userID)
	order := o.OrderRepository.FindByID(orderID)
	if user.Customer.ID != order.Customer.ID {
		err := errors.New("订单 ID 与用户 ID 不符合!")
		o.Logger.Log(err.Error())
		return err
	}
	item := o.ItemRepository.FindByID(itemID)
	if domainErr := order.Add(item); domainErr != nil {
		o.Logger.Log(domainErr.Error())
		return domainErr
	}

	o.OrderRepository.Store(order)
	return nil
}

type AdminOrderInteractor struct {
	OrderInteractor
}

func (a *AdminOrderInteractor) Add(userID, orderID, itemID int) error {
	user := a.UserRepository.FindByID(userID)
	order := a.OrderRepository.FindByID(orderID)

	if !user.IsAdmin {
		err := errors.New("用户不是 admin!")
		a.Logger.Log(err.Error())
		return err
	}

	item := a.ItemRepository.FindByID(itemID)
	if domainErr := order.Add(item); domainErr != nil {
		a.Logger.Log(domainErr.Error())
		return domainErr
	}
	a.OrderRepository.Store(order)
	return nil
}
