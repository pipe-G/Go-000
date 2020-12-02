package Week02

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

const(
	USERNAME="root"
	PASSWORD="123456"
	NETWORK="tcp"
	SERVER="127.0.0.1"
	PORT="3306"
	DATABASE="test"
)

type Order struct{
	OrderID    int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func GetOrder(OrderID string) (*Order, error) {
	connection:=fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME,PASSWORD,NETWORK,SERVER,PORT,DATABASE)
	db,err:=sql.Open("mysql",connection)
	if err !=nil{
		return nil, errors.Wrapf(err,"connection failure")
	}
	defer  db.Close()
	row :=db.QueryRow("select * from order where OrderID=?",OrderID)
	order :=new(Order)
	if err := row.Scan(&order.OrderID,&order.CreatedAt,&order.UpdatedAt); err != nil {
		return nil, errors.Wrapf(err,"data not found")
	}
	return order, nil
}

func OrderService()(*Order,error) {
	return GetOrder("1111111111")
}

func main() {
	order,err:=OrderService()
	if err!=nil{
		fmt.Println("error msg %v",err)
	}
	fmt.Println("create time is %s",order.CreatedAt)
}

