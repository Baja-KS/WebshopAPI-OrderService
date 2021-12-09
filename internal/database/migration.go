package database

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Order{},&OrderItem{})
	if err != nil {
		return err
	}
	return nil
}
