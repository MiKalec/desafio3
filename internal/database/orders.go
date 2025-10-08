package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Order struct {
	db          *sql.DB
	ID          string
	Name        string
	OrderType   string
	Description string
}

func NewOrder(db *sql.DB) *Order {
	return &Order{db: db}
}

func (o *Order) Create(name string, orderType string, description string) (Order, error) {
	id := uuid.New().String()
	_, err := o.db.Exec("insert into orders (id, name, orderType, description) VALUES (?, ?, ?, ?)", id, name, orderType, description)
	if err != nil {
		return Order{}, err
	}
	return Order{ID: id, Name: name, OrderType: orderType, Description: description}, nil
}

func (o *Order) FindAll() ([]Order, error) {
	rows, err := o.db.Query("SELECT id, name, orderType, description FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var id, name, orderType, description string
		if err := rows.Scan(&id, &name, &orderType, &description); err != nil {
			return nil, err
		}
		orders = append(orders, Order{ID: id, Name: name, OrderType: orderType, Description: description})
	}
	return orders, nil
}
