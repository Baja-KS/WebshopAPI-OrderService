package database

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

type OrderItem struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Count uint `gorm:"default:1" json:"count"`
	OrderID uint `gorm:"not null" json:"OrderId"`
	ProductID uint `gorm:"not null" json:"ProductId"`
}

type OrderItemIn struct {
	Count uint `json:"count"`
	ProductID uint `json:"ProductId"`
}

type OrderItemOut struct {
	ID uint `json:"id,omitempty"`
	Count uint `json:"count"`
	OrderID uint `json:"OrderId"`
	ProductID uint `json:"ProductId"`
	Product ProductOut `json:"Product"`
}


func (i *OrderItem) GetProduct(productServiceUrl string) (ProductOut,error) {
	var product ProductOut
	var response ProductServiceResponse
	res,err:=http.Get(productServiceUrl+"/GetByID/"+strconv.Itoa(int(i.ProductID)))
	if err != nil {
		return product,err
	}
	err=json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return product,err
	}
	return response.Product,nil
}

func (i *OrderItem) GetValue() float32 {
	product,err:=i.GetProduct(os.Getenv("PRODUCT_SERVICE"))
	if err != nil {
		return 0
	}
	return product.Price*float32(i.Count)
}

func (i *OrderItem) Out() OrderItemOut {
	product,_:=i.GetProduct(os.Getenv("PRODUCT_SERVICE"))
	return OrderItemOut{
		ID:        i.ID,
		Count:     i.Count,
		OrderID:   i.OrderID,
		ProductID: i.ProductID,
		Product: product,
	}
}

func (i *OrderItemIn) In() OrderItem {
	return OrderItem{
		Count:     i.Count,
		ProductID: i.ProductID,
	}
}

func ItemArrayOut(models []OrderItem) []OrderItemOut {
	outArr:=make([]OrderItemOut,len(models))
	for i,item := range models {
		outArr[i]=item.Out()
	}
	return outArr
}

func ItemArrayIn(models []OrderItemIn) []OrderItem {
	inArr:=make([]OrderItem,len(models))
	for i, item := range models {
		inArr[i]=item.In()
	}
	return inArr
}

