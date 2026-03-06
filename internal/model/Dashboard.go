package model

import "time"

type CategorySpending struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type MonthlySpending struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

type SupplierSpending struct {
	SupplierName string  `json:"supplier_name"`
	Total        float64 `json:"total"`
}

type DashboardStats struct {
	TotalSessions      int64  `json:"total_sessions"`
	TotalProducts      int64  `json:"total_products"`
	MostBoughtCategory string `json:"most_bought_category"`
	MostBoughtProduct  string `json:"most_bought_product"`
}

type DashboardSummary struct {
	Stats          *DashboardStats    `json:"stats"`
	ByCategory     []CategorySpending `json:"by_category"`
	ByMonth        []MonthlySpending  `json:"by_month"`
	BySupplier     []SupplierSpending `json:"by_supplier"`
	RecentSessions []PurchaseSession  `json:"recent_sessions"`
}

type PriceHistory struct {
	Product       string    `json:"product"`
	LatestPrice   float64   `json:"latest_price"`
	PreviousPrice float64   `json:"previous_price"`
	Change        float64   `json:"change"`
	LastPurchased time.Time `json:"last_purchased"`
}

type ReorderReminder struct {
    Product               string    `json:"product"`
    Category              string    `json:"category"`
    LastPurchased         time.Time `json:"last_purchased"`
    DaysSinceLastPurchase int       `json:"days_since_last_purchase"`
}

//admin dashboard
type AdminStats struct {
	TotalUsers             int64 `json:"total_users"`
	TotalSessions          int64 `json:"total_sessions"`
	TotalProductItems      int64 `json:"total_product_items"`
	ActiveUsersLast7Days   int64 `json:"active_users_last_7_days"`
	ActiveUsersLast30Days  int64 `json:"active_users_last_30_days"`
	NewUsersThisMonth      int64 `json:"new_users_this_month"`
}

type AdminDashboardSummary struct {
	Stats *AdminStats `json:"stats"`
}

type PaginatedSessions struct {
	Sessions []PurchaseSession `json:"sessions"`
	Meta     PaginationMeta    `json:"meta"`
}

type PaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"has_more"`
}