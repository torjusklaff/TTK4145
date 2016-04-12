package queue

import (
	"../driver"
	"../network"
	"fmt"
)

type OrderButton struct {
	Type	driver.ButtonType
	Floor	int
}

type Order struct {
	Button	driver.OrderButton
	TakenBy	network.IP
	Done	bool
}

func OrdersEqual(order1, order2 Order) bool {
	return	order1.Button.Floor == order2.Button.Floor &&
			order1.Button.Type == order2.Button.Type
}

func OrderNew(request Order, orders []Order) bool {
	for _, order := range(orders) {
		if OrdersEqual(request, order) {
			return false
		}
	}
	return true
}

func GetPriority(orders []Order, ip network.IP) *Order {
	for _, order := range(orders) {
		if order.TakenBy == id && order.Priority {
			return &order
		}
	}
	return nil
}
