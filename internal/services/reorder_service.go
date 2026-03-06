package services

// import (
// 	"context"

// 	"github.com/google/uuid"
// 	"github.com/britinogn/bizkeeper/internal/model"
// )

// type ReorderRepo interface {
// 	GetReorderReminders(ctx context.Context, userID uuid.UUID, dayThreshold int) ([]model.ReorderReminder, error)
// }

// type ReorderService struct {
// 	reorderRepo ReorderRepo
// }

// func NewReorderService(reorderRepo ReorderRepo) *ReorderService {
// 	return &ReorderService{reorderRepo: reorderRepo}
// }

// func (s *ReorderService) GetReorderReminders(ctx context.Context, userID uuid.UUID) ([]model.ReorderReminder, error) {
// 	const dayThreshold = 14
// 	return s.reorderRepo.GetReorderReminders(ctx, userID, dayThreshold)
// }