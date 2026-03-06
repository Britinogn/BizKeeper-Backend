package repository

import (
	"context"
	"fmt"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CreatePurchaseSession(ctx context.Context, session *model.PurchaseSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *PurchaseRepository) GetPurchaseSessionByID(ctx context.Context, id uuid.UUID) (*model.PurchaseSession, error) {
	var session model.PurchaseSession
	// Eager load product items
	err := r.db.WithContext(ctx).Preload("ProductItems").Where("id = ?", id).First(&session).Error
	return &session, err
}

func (r *PurchaseRepository) ListPurchaseSessions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.PurchaseSession, error) {
	var sessions []model.PurchaseSession
	// Eager load product items for each session
	err := r.db.WithContext(ctx).Preload("ProductItems").Where("user_id = ?", userID).Order("purchase_date desc").Limit(limit).Offset(offset).Find(&sessions).Error
	return sessions, err
}

func (r *PurchaseRepository) CountPurchaseSessions(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PurchaseSession{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *PurchaseRepository) UpdatePurchaseSession(ctx context.Context, session *model.PurchaseSession) error {
	// Start a transaction to handle both session and items
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update the session
		if err := tx.Save(session).Error; err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}

		// Delete existing items for this session
		if err := tx.Where("session_id = ?", session.ID).Delete(&model.ProductItem{}).Error; err != nil {
			return fmt.Errorf("failed to delete old items: %w", err)
		}

		// Insert new items
		if len(session.ProductItems) > 0 {
			// Set session_id for all items
			for i := range session.ProductItems {
				session.ProductItems[i].SessionID = session.ID
			}
			if err := tx.Create(&session.ProductItems).Error; err != nil {
				return fmt.Errorf("failed to insert new items: %w", err)
			}
		}

		return nil
	})
}

func (r *PurchaseRepository) DeletePurchaseSession(ctx context.Context, session *model.PurchaseSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hard delete items first using raw SQL
		if err := tx.Exec("DELETE FROM product_items WHERE session_id = ?", session.ID).Error; err != nil {
			return fmt.Errorf("failed to delete items: %w", err)
		}

		// Hard delete the session
		if err := tx.Unscoped().Delete(session).Error; err != nil {
			return fmt.Errorf("failed to delete session: %w", err)
		}

		return nil
	})
}


// BELOW IS TO FETCH A SINGLE Item SESSION

func (r *PurchaseRepository) UpdateProductItem(ctx context.Context, item *model.ProductItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *PurchaseRepository) DeleteProductItem(ctx context.Context, item *model.ProductItem) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM product_items WHERE id = ?", item.ID).Error
}

func (r *PurchaseRepository) GetProductItemByID(ctx context.Context, itemID uuid.UUID) (*model.ProductItem, error) {
	var item model.ProductItem
	return &item, r.db.WithContext(ctx).Where("id = ?", itemID).First(&item).Error
}