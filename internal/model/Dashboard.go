package model

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