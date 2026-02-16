package usecase

import (
	"testing"

	"github.com/MiKalec/desafio3/internal/entity"
	"github.com/stretchr/testify/assert"
)

type OrderRepositoryListMock struct {
	Orders []entity.Order
	Err    error
}

func (m *OrderRepositoryListMock) Save(order *entity.Order) error {
	return nil
}

func (m *OrderRepositoryListMock) GetAll() ([]entity.Order, error) {
	return m.Orders, m.Err
}

func TestListOrdersUseCase_Execute(t *testing.T) {
	orders := []entity.Order{
		{
			ID:         "123",
			Price:      10.0,
			Tax:        2.0,
			FinalPrice: 12.0,
		},
		{
			ID:         "456",
			Price:      20.0,
			Tax:        3.0,
			FinalPrice: 23.0,
		},
	}

	repo := &OrderRepositoryListMock{
		Orders: orders,
	}

	usecase := NewListOrdersUseCase(repo)
	output, err := usecase.Execute()

	assert.NoError(t, err)
	assert.Len(t, output, 2)

	assert.Equal(t, "123", output[0].ID)
	assert.Equal(t, 10.0, output[0].Price)
	assert.Equal(t, 2.0, output[0].Tax)
	assert.Equal(t, 12.0, output[0].FinalPrice)

	assert.Equal(t, "456", output[1].ID)
	assert.Equal(t, 20.0, output[1].Price)
	assert.Equal(t, 3.0, output[1].Tax)
	assert.Equal(t, 23.0, output[1].FinalPrice)
}
