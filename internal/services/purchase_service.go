package services

import (
	"context"
	"errors"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound     = errors.New("purchase session not found")
	ErrUnauthorizedSession = errors.New("unauthorized to access this session")
	ErrInvalidSession      = errors.New("invalid session data")
	ErrNoProductItems      = errors.New("at least one product item is required")
	ErrItemNotFound        = errors.New("product item not found")
	//ErrDatabaseOperation   = errors.New("database operation failed")
)

type PurchaseRepo interface {
	CreatePurchaseSession(ctx context.Context, session *model.PurchaseSession) error
	GetPurchaseSessionByID(ctx context.Context, id uuid.UUID) (*model.PurchaseSession, error)
	ListPurchaseSessions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.PurchaseSession, error)
	UpdatePurchaseSession(ctx context.Context, session *model.PurchaseSession) error
	DeletePurchaseSession(ctx context.Context, session *model.PurchaseSession) error

	UpdateProductItem(ctx context.Context, item *model.ProductItem) error
	DeleteProductItem(ctx context.Context, item *model.ProductItem) error
	GetProductItemByID(ctx context.Context, itemID uuid.UUID) (*model.ProductItem, error)
}

type PurchaseService struct {
	purchaseRepo PurchaseRepo
}

func NewPurchaseService(purchaseRepo PurchaseRepo) *PurchaseService {
	return &PurchaseService{purchaseRepo: purchaseRepo}
}

func (s *PurchaseService) CreatePurchaseSession(ctx context.Context, userID uuid.UUID, session *model.PurchaseSession) error {
	if session == nil {
		return ErrInvalidSession
	}

	if session.SupplierName == "" {
		return ErrInvalidSession
	}

	if len(session.ProductItems) == 0 {
		return ErrNoProductItems
	}

	// Always tie session to authenticated user
	session.UserID = userID

	// Set session ID on all product items
	// for i := range session.ProductItems {
	// 	session.ProductItems[i].SessionID = session.ID
	// }
	
	for i := range session.ProductItems {
		session.ProductItems[i].SubtotalAmount = float64(session.ProductItems[i].Quantity) * session.ProductItems[i].UnitPrice
		session.TotalAmount += session.ProductItems[i].SubtotalAmount
	}

	return s.purchaseRepo.CreatePurchaseSession(ctx, session)
}

func (s *PurchaseService) GetPurchaseSessionByID(ctx context.Context, userID, sessionID uuid.UUID) (*model.PurchaseSession, error) {
	session, err := s.purchaseRepo.GetPurchaseSessionByID(ctx, sessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Ensure session belongs to the authenticated user
	if session.UserID != userID {
		return nil, ErrUnauthorizedSession
	}

	// Calculate subtotals and total
	for i := range session.ProductItems {
		session.ProductItems[i].SubtotalAmount = float64(session.ProductItems[i].Quantity) * session.ProductItems[i].UnitPrice
		session.TotalAmount += session.ProductItems[i].SubtotalAmount
	}


	return session, nil
}

func (s *PurchaseService) ListPurchaseSessions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.PurchaseSession, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	sessions, err := s.purchaseRepo.ListPurchaseSessions(ctx, userID, limit, offset)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	for i := range sessions {
    for j := range sessions[i].ProductItems {
        sessions[i].ProductItems[j].SubtotalAmount = float64(sessions[i].ProductItems[j].Quantity) * sessions[i].ProductItems[j].UnitPrice
        sessions[i].TotalAmount += sessions[i].ProductItems[j].SubtotalAmount
    }
}

	return sessions, nil
}

func (s *PurchaseService) UpdatePurchaseSession(ctx context.Context, userID, sessionID uuid.UUID, updated *model.PurchaseSession) error {
	session, err := s.purchaseRepo.GetPurchaseSessionByID(ctx, sessionID)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.UserID != userID {
		return ErrUnauthorizedSession
	}

	updated.ID = session.ID
	updated.UserID = userID

	return s.purchaseRepo.UpdatePurchaseSession(ctx, updated)
}

func (s *PurchaseService) DeletePurchaseSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	session, err := s.purchaseRepo.GetPurchaseSessionByID(ctx, sessionID)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.UserID != userID {
		return ErrUnauthorizedSession
	}

	return s.purchaseRepo.DeletePurchaseSession(ctx, session)
}




func (s *PurchaseService) UpdateProductItem(ctx context.Context, userID, sessionID, itemID uuid.UUID, updated *model.ProductItem) (*model.ProductItem, error) {
	// Verify session exists and belongs to user
	session, err := s.purchaseRepo.GetPurchaseSessionByID(ctx, sessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	if session.UserID != userID {
		return nil, ErrUnauthorizedSession
	}

	// Verify item exists and belongs to session
	item, err := s.purchaseRepo.GetProductItemByID(ctx, itemID)
	if err != nil {
		return nil, ErrItemNotFound
	}
	if item.SessionID != sessionID {
		return nil, ErrUnauthorizedSession
	}

	// Update only provided fields
	if updated.Name != "" {
		item.Name = updated.Name
	}
	if updated.Quantity > 0 {
		item.Quantity = updated.Quantity
	}
	if updated.UnitPrice > 0 {
		item.UnitPrice = updated.UnitPrice
	}
	if updated.Category != "" {
		item.Category = updated.Category
	}
	if updated.Notes != "" {
		item.Notes = updated.Notes
	}

	if err := s.purchaseRepo.UpdateProductItem(ctx, item); err != nil {
		return nil, ErrDatabaseOperation
	}

	return item, nil
}

func (s *PurchaseService) DeleteProductItem(ctx context.Context, userID, sessionID, itemID uuid.UUID) error {
	// Verify session exists and belongs to user
	session, err := s.purchaseRepo.GetPurchaseSessionByID(ctx, sessionID)
	if err != nil {
		return ErrSessionNotFound
	}
	if session.UserID != userID {
		return ErrUnauthorizedSession
	}

	// Verify item exists and belongs to session
	item, err := s.purchaseRepo.GetProductItemByID(ctx, itemID)
	if err != nil {
		return ErrItemNotFound
	}
	if item.SessionID != sessionID {
		return ErrUnauthorizedSession
	}

	return s.purchaseRepo.DeleteProductItem(ctx, item)
}