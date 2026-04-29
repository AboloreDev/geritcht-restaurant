package models

import "time"

type DailySummary struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Date           time.Time `json:"date" gorm:"not null;uniqueIndex"`
	TotalOrders    int       `json:"total_orders" gorm:"default:0"`
	TotalRevenue   float64   `json:"total_revenue" gorm:"default:0"`
	TotalCustomers int       `json:"total_customers" gorm:"default:0"`
	PopularItemID  *uint     `json:"popular_item_id"`
	CreatedAt      time.Time `json:"created_at"`
	// Relationships
	PopularItem *Menu `json:"popular_item,omitempty"`
}
