package database

import (
	"database/sql"
	"testing"

	"github.com/MiKalec/desafio3/internal/entity"
	"github.com/stretchr/testify/suite"

	_ "modernc.org/sqlite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("sqlite", ":memory:")
	suite.NoError(err)
	_, err = db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	if suite.Db != nil {
		suite.Db.Close()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestGivenOrders_WhenGetAll_ThenShouldReturnAllOrders() {
	repo := NewOrderRepository(suite.Db)

	order1, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order1.CalculateFinalPrice())

	order2, err := entity.NewOrder("456", 20.0, 3.0)
	suite.NoError(err)
	suite.NoError(order2.CalculateFinalPrice())

	err = repo.Save(order1)
	suite.NoError(err)
	err = repo.Save(order2)
	suite.NoError(err)

	orders, err := repo.GetAll()
	suite.NoError(err)
	suite.Len(orders, 2)

	foundIDs := make(map[string]bool)
	for _, o := range orders {
		foundIDs[o.ID] = true
	}
	suite.True(foundIDs[order1.ID])
	suite.True(foundIDs[order2.ID])
}
