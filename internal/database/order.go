package database

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	FirstName string `gorm:"not null" json:"firstName"`
	LastName string `gorm:"not null" json:"lastName"`
	Email string `gorm:"not null" json:"email"`
	Address string `gorm:"not null" json:"address"`
	City string `gorm:"not null" json:"city"`
	PaymentMethod string `gorm:"not null" json:"paymentMethod"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
	OrderItems []OrderItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"orderItems"`
}

type OrderIn struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Address string `json:"address"`
	City string `json:"city"`
	PaymentMethod string `json:"paymentMethod"`
	OrderItems []OrderItemIn `json:"items"`
}

type OrderOut struct {
	ID uint `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Address string `json:"address"`
	City string `json:"city"`
	PaymentMethod string `json:"paymentMethod"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Total float32 `json:"total"`
}

func (o *Order) GetTotalValue() float32 {
	var total float32
	total=0
	for _, item := range o.OrderItems {
		total+=item.GetValue()
	}
	return total
}

func GetTotalValue(orders []Order) float32 {
	var total float32
	total=0
	for _, order := range orders {
		total+=order.GetTotalValue()
	}
	return total
}

func (o *Order) Out() OrderOut {
	return OrderOut{
		ID:            o.ID,
		FirstName:     o.FirstName,
		LastName:      o.LastName,
		Email:         o.Email,
		Address:       o.Address,
		City:          o.City,
		PaymentMethod: o.PaymentMethod,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Total: o.GetTotalValue(),
	}
}

func (i *OrderIn) In() Order {
	return Order{
		FirstName:     i.FirstName,
		LastName:      i.LastName,
		Email:         i.Email,
		Address:       i.Address,
		City:          i.City,
		PaymentMethod: i.PaymentMethod,
		OrderItems:    ItemArrayIn(i.OrderItems),
	}
}

func OrderArrayOut(models []Order) []OrderOut {
	outArr:=make([]OrderOut,len(models))
	for i,item := range models {
		outArr[i]=item.Out()
	}
	return outArr
}


