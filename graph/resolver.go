package graph

import "github.com/MiKalec/desafio3/internal/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	OrderDB *database.Order
}
