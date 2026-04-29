package dto

import "time"

type AnalyticsFilterRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Period    string `form:"period,default=daily"`
}

type DailySummaryResponse struct {
	Date           string        `json:"date"`
	TotalOrders    int           `json:"total_orders"`
	TotalRevenue   float64       `json:"total_revenue"`
	TotalCustomers int           `json:"total_customers"`
	PopularItem    *MenuResponse `json:"popular_item,omitempty"`
}

type RevenueResponse struct {
	Period       string    `json:"period"`
	TotalRevenue float64   `json:"total_revenue"`
	TotalOrders  int       `json:"total_orders"`
	AverageOrder float64   `json:"average_order_value"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type PeakHoursResponse struct {
	Hour        int `json:"hour"`
	TotalOrders int `json:"total_orders"`
}

type PopularItemResponse struct {
	MenuItem    MenuResponse `json:"menu_item"`
	TotalOrders int          `json:"total_orders"`
	Revenue     float64      `json:"revenue"`
}

type AnalyticsSummaryResponse struct {
	Revenue          RevenueResponse        `json:"revenue"`
	PeakHours        []PeakHoursResponse    `json:"peak_hours"`
	PopularItems     []PopularItemResponse  `json:"popular_items"`
	TableUtilization float64                `json:"table_utilization_rate"`
	DailySummaries   []DailySummaryResponse `json:"daily_summaries"`
}
