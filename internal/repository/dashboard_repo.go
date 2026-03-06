package repository

import (
	"context"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/google/uuid"
)

// type PurchaseRepository struct {
// 	db *gorm.DB
// }

// func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
// 	return &PurchaseRepository{db: db}
// }

// type DashboardRepository struct {
// 	db *gorm.DB
// }

// func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
// 	return &DashboardRepository{db: db}
// }

func (r *PurchaseRepository) GetSpendingByCategory(ctx context.Context, userID uuid.UUID) ([]model.CategorySpending, error) {
	var results []model.CategorySpending
	err := r.db.WithContext(ctx).
		Table("product_items pi").
		Select("pi.category, SUM(pi.quantity * pi.unit_price) as total").
		Joins("JOIN purchase_sessions ps ON ps.id = pi.session_id").
		Where("ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL", userID).
		Group("pi.category").
		Scan(&results).Error
	return results, err
}

func (r *PurchaseRepository) GetSpendingByMonth(ctx context.Context, userID uuid.UUID) ([]model.MonthlySpending, error) {
	var results []model.MonthlySpending
	err := r.db.WithContext(ctx).
		Table("product_items pi").
		Select("TO_CHAR(ps.purchase_date, 'YYYY-MM') as month, SUM(pi.quantity * pi.unit_price) as total").
		Joins("JOIN purchase_sessions ps ON ps.id = pi.session_id").
		Where("ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL", userID).
		Group("month").
		Order("month desc").
		Scan(&results).Error
	return results, err
}

func (r *PurchaseRepository) GetSpendingBySupplier(ctx context.Context, userID uuid.UUID) ([]model.SupplierSpending, error) {
	var results []model.SupplierSpending
	err := r.db.WithContext(ctx).
		Table("product_items pi").
		Select("ps.supplier_name, SUM(pi.quantity * pi.unit_price) as total").
		Joins("JOIN purchase_sessions ps ON ps.id = pi.session_id").
		Where("ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL", userID).
		Group("ps.supplier_name").
		Scan(&results).Error
	return results, err
}

func (r *PurchaseRepository) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*model.DashboardStats, error) {
	var stats model.DashboardStats

	// Total sessions
	r.db.WithContext(ctx).Model(&model.PurchaseSession{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&stats.TotalSessions)

	// Total product items
	r.db.WithContext(ctx).
		Table("product_items pi").
		Joins("JOIN purchase_sessions ps ON ps.id = pi.session_id").
		Where("ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL", userID).
		Count(&stats.TotalProducts)

	// Most bought category
	r.db.WithContext(ctx).
		Table("product_items pi").
		Select("pi.category").
		Joins("JOIN purchase_sessions ps ON ps.id = pi.session_id").
		Where("ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL", userID).
		Group("pi.category").
		Order("SUM(pi.quantity) desc").
		Limit(1).
		Pluck("pi.category", &stats.MostBoughtCategory)

	// Most bought product
	r.db.WithContext(ctx).
		Table("product_items pi").
		Select("pi.name").
		Joins("JOIN purchase_sessions ps ON ps.id = pi.session_id").
		Where("ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL", userID).
		Group("pi.name").
		Order("SUM(pi.quantity) desc").
		Limit(1).
		Pluck("pi.name", &stats.MostBoughtProduct)

	return &stats, nil
}

func (r *PurchaseRepository) GetRecentSessions(ctx context.Context, userID uuid.UUID) ([]model.PurchaseSession, error) {
	var sessions []model.PurchaseSession
	err := r.db.WithContext(ctx).
		Preload("ProductItems").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("purchase_date desc").
		Limit(5).
		Find(&sessions).Error
	return sessions, err
}


func (r *PurchaseRepository) GetPriceHistory(ctx context.Context, userID uuid.UUID) ([]model.PriceHistory, error) {
	var results []model.PriceHistory

	err := r.db.WithContext(ctx).Raw(`
		WITH ranked_prices AS (
			SELECT 
				pi.name,
				pi.unit_price,
				ps.purchase_date,
				ROW_NUMBER() OVER (PARTITION BY pi.name ORDER BY ps.purchase_date DESC) as rn
			FROM product_items pi
			JOIN purchase_sessions ps ON ps.id = pi.session_id
			WHERE ps.user_id = ? AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL
		)
		SELECT
			a.name as product,
			a.unit_price as latest_price,
			b.unit_price as previous_price,
			a.unit_price - COALESCE(b.unit_price, a.unit_price) as change,
			a.purchase_date as last_purchased
		FROM ranked_prices a
		LEFT JOIN ranked_prices b ON a.name = b.name AND b.rn = 2
		WHERE a.rn = 1
		ORDER BY a.name
	`, userID).Scan(&results).Error

	return results, err
}