package services

import (
	"context"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/google/uuid"
)

type DashboardRepo interface {
	GetSpendingByCategory(ctx context.Context, userID uuid.UUID) ([]model.CategorySpending, error)
	GetSpendingByMonth(ctx context.Context, userID uuid.UUID) ([]model.MonthlySpending, error)
	GetSpendingBySupplier(ctx context.Context, userID uuid.UUID) ([]model.SupplierSpending, error)
	GetDashboardStats(ctx context.Context, userID uuid.UUID) (*model.DashboardStats, error)
	GetRecentSessions(ctx context.Context, userID uuid.UUID) ([]model.PurchaseSession, error)
	GetPriceHistory(ctx context.Context, userID uuid.UUID) ([]model.PriceHistory, error)
	GetAdminDashboardStats(ctx context.Context) (*model.AdminStats, error)
	GetReorderReminders(ctx context.Context, userID uuid.UUID) ([]model.ReorderReminder, error)
}

type DashboardService struct {
	dashboardRepo DashboardRepo
}

func NewDashboardService(dashboardRepo DashboardRepo) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
}

func (s *DashboardService) GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*model.DashboardSummary, error) {
	// Stats
	stats, err := s.dashboardRepo.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Spending by category
	byCategory, err := s.dashboardRepo.GetSpendingByCategory(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Spending by month
	byMonth, err := s.dashboardRepo.GetSpendingByMonth(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Spending by supplier
	bySupplier, err := s.dashboardRepo.GetSpendingBySupplier(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Recent sessions
	recentSessions, err := s.dashboardRepo.GetRecentSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate totals for recent sessions
	for i := range recentSessions {
		for j := range recentSessions[i].ProductItems {
			recentSessions[i].ProductItems[j].SubtotalAmount = float64(recentSessions[i].ProductItems[j].Quantity) * recentSessions[i].ProductItems[j].UnitPrice
			recentSessions[i].TotalAmount += recentSessions[i].ProductItems[j].SubtotalAmount
		}
	}

	// calculate total spend from category totals
	var totalSpend float64
	for _, c := range byCategory {
		totalSpend += c.Total
	}
	stats.TotalSpend = totalSpend

	return &model.DashboardSummary{
		Stats:          stats,
		ByCategory:     byCategory,
		ByMonth:        byMonth,
		BySupplier:     bySupplier,
		RecentSessions: recentSessions,
	}, nil
}

func (s *DashboardService) GetPriceHistory(ctx context.Context, userID uuid.UUID) ([]model.PriceHistory, error) {
	return s.dashboardRepo.GetPriceHistory(ctx, userID)
}

// admin dashboard only
func (s *DashboardService) GetAdminDashboard(ctx context.Context) (*model.AdminStats, error) {
	return s.dashboardRepo.GetAdminDashboardStats(ctx)
}


func (s *DashboardService) GetReorderReminders(ctx context.Context, userID uuid.UUID) ([]model.ReorderReminder, error) {
	return s.dashboardRepo.GetReorderReminders(ctx, userID)
}