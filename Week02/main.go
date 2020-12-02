package Week02

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
	goqu "gopkg.in/doug-martin/goqu.v5"
)

type Order struct{
	OrderID    int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (t Order) TableName() string {
	return "Order"
}


func OrderBuilder(params map[string]interface{}) *goqu.Dataset {
	return db.DB.Goqu.From(goqu.I(Order{}.TableName()).As("o")).
		Where(
			db.QueryFilter(params).
				Where("o.order_id", db.EQ, "OrderID").
				End()...,
		).Order(goqu.I("o.created_at").Desc())
}

func GetOrder(params map[string]interface{}) (*Order, error) {
	var (
		order Order
	)
	builder := OrderBuilder(params)
	if err := db.DB.QueryFirst(builder, &order); err != nil {
		return nil, errors.Wrapf(err,"data not found")
	}
	return &order, nil
}

func OrderService()(*Order,error) {
	return GetOrder(map[string]interface{}{
		"OutOrderID": "11111111",
	})
}

func main() {
	order,err:=OrderService()
	if err!=nil{
		fmt.Println("error msg %v",err)
	}
	fmt.Println("create time is %s",order.CreatedAt)
}

