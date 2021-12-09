package service

import (
	"OrderService/internal/database"
	"context"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

//OrderService should implement the Service interface


type OrderService struct {
	DB *gorm.DB
}

func ValidateProduct(productServiceURL string, id uint) bool {
	_,err:=http.Get(productServiceURL+"/GetByID/"+ strconv.Itoa(int(id)))
	if err != nil {
		return false
	}
	return true
}

type Service interface {
	GetByID(ctx context.Context, id uint) ([]database.OrderItemOut, error)
	Search(ctx context.Context,search string,startDate time.Time,endDate time.Time) ([]database.OrderOut,error)
	Create(ctx context.Context,data database.OrderIn) (string,error)
	Delete(ctx context.Context,id uint) (string,error)
	Total(ctx context.Context) (float32,error)
	Top(ctx context.Context, count uint) ([]database.ProductOutTop, error)
	QuantityOrdered(ctx context.Context, id uint) (uint, error)
}

func (o *OrderService) QuantityOrdered(ctx context.Context, id uint) (uint, error) {
	var items []database.OrderItem
	result:=o.DB.Where("product_id = ?",id).Find(&items)
	var qty uint
	qty=0
	if result.Error != nil {
		return qty,nil
	}
	for _, item := range items {
		qty+=item.Count
	}
	return qty,nil
}

func (o *OrderService) GetByID(ctx context.Context, id uint) ([]database.OrderItemOut, error) {
	var items []database.OrderItem
	result:=o.DB.Where("order_id = ?",id).Find(&items)
	if result.Error != nil {
		return nil,result.Error
	}
	return database.ItemArrayOut(items),nil
}

func (o *OrderService) Search(ctx context.Context, search string, startDate time.Time, endDate time.Time) ([]database.OrderOut, error) {
	var orders []database.Order
	result:=o.DB.Where(
		o.DB.Where("cast(id as varchar) ilike ?","%"+search+"%").Or("first_name ilike ?","%"+search+"%").Or("last_name ilike ?","%"+search+"%").Or("email ilike ?","%"+search+"%").Or("address ilike ?","%"+search+"%").Or("city ilike ?","%"+search+"%"),
	)
	if !startDate.IsZero() && !startDate.Equal(time.Unix(0,0)) {
		result=result.Where("created_at >= ?",startDate)
	}
	if !endDate.IsZero() && !endDate.Equal(time.Unix(0,0)) {
		result=result.Where("created_at <= ?",endDate)
	}
	if result.Preload("OrderItems").Find(&orders).Error != nil {
		return database.OrderArrayOut(orders),result.Error
	}
	return database.OrderArrayOut(orders),nil
}

func (o *OrderService) Create(ctx context.Context, data database.OrderIn) (string, error) {
	order:=data.In()
	for _, item := range order.OrderItems {
		if !ValidateProduct(os.Getenv("PRODUCT_SERVICE"), item.ProductID) {
			return "Product with an id "+ strconv.Itoa(int(item.ProductID)) +" doesn't exist", errors.New("product with that ID doesnt exist")
		}
	}

	result:=o.DB.Create(&order)
	if result.Error != nil {
		return "Error", result.Error
	}
	return "Order successful", nil

}

func (o *OrderService) Delete(ctx context.Context, id uint) (string, error) {
	var order database.Order
	notFound:=o.DB.Where("id = ?",id).First(&order).Error
	if notFound != nil {
		return "That order doesn't exist", notFound
	}
	err:=o.DB.Delete(&database.Order{},id).Error
	if err != nil {
		return "Error deleting order", err
	}

	return "Order deleted successfully", nil
}

func (o *OrderService) Total(ctx context.Context) (float32, error) {
	var orders []database.Order
	result:=o.DB.Preload("OrderItems").Find(&orders)
	if result.Error != nil {
		return 0, nil
	}
	return database.GetTotalValue(orders),nil
}

func (o *OrderService) Top(ctx context.Context, count uint) ([]database.ProductOutTop, error) {
	var items []database.OrderItem
	result:=o.DB.Find(&items)
	if result.Error != nil {
		return nil,result.Error
	}
	var products []database.ProductOutTop
	for _, item := range items {
		product,err:=item.GetProduct(os.Getenv("PRODUCT_SERVICE"))
		if err == nil {
			products=append(products,product.Top(product.ID,o.DB))
		}
	}
	sort.Sort(sort.Reverse(database.TopSellingProducts(products)))

	n:=uint(len(products))

	var elCount uint

	if n < count {
		elCount = n
	} else {
		elCount = count
	}

	products=products[0:elCount]

	return products,nil
}
