package database

import "gorm.io/gorm"

type ProductOut struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Img string `json:"img,omitempty"`
	Price float32 `json:"price"`
	Discount int `json:"discount"`
	CategoryID uint `json:"CategoryId"`
	Deletable bool `json:"deletable"`
}

func (p *ProductOut) Top(ProductID uint,db *gorm.DB) ProductOutTop {
	return ProductOutTop{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Img:         p.Img,
		Price:       p.Price,
		Discount:    p.Discount,
		CategoryID:  p.CategoryID,
		Deletable: p.Deletable,
		TimesSold:   timesSold(ProductID,db),
	}
}

func timesSold(productID uint,db *gorm.DB) uint {
	var items []OrderItem
	var timesSold uint
	timesSold=0
	result:=db.Where("product_id = ?",productID).Find(&items)
	if result.Error != nil {
		return timesSold
	}
	for _, item := range items {
		timesSold+=item.Count
	}
	return timesSold
}

type ProductServiceResponse struct {
	Product ProductOut `json:"product"`
}

type ProductOutTop struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Img string `json:"img,omitempty"`
	Price float32 `json:"price"`
	Discount int `json:"discount"`
	CategoryID uint `json:"CategoryId"`
	Deletable bool `json:"deletable"`
	TimesSold uint `json:"timesSold"`
}

type TopSellingProducts []ProductOutTop

func (t TopSellingProducts) Len() int {
	return len(t)
}

func (t TopSellingProducts) Less(i, j int) bool {
	return t[i].TimesSold < t[j].TimesSold
}

func (t TopSellingProducts) Swap(i, j int) {
	t[i],t[j]=t[j],t[i]
}
