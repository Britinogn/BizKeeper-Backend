package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type ExportRepo interface {
	GetSessionsByDateRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.PurchaseSession, error)
}

type ExportService struct {
	exportRepo ExportRepo
}

func NewExportService(exportRepo ExportRepo) *ExportService {
	return &ExportService{exportRepo: exportRepo}
}

func parseRange(rangeStr string) (time.Time, time.Time, error) {
	to := time.Now()
	var from time.Time

	switch rangeStr {
	case "7days":
		from = to.AddDate(0, 0, -7)
	case "1month":
		from = to.AddDate(0, -1, 0)
	case "3months":
		from = to.AddDate(0, -3, 0)
	case "6months":
		from = to.AddDate(0, -6, 0)
	case "1year":
		from = to.AddDate(-1, 0, 0)
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid range: %s", rangeStr)
	}

	return from, to, nil
}

func (s *ExportService) ExportCSV(ctx context.Context, userID uuid.UUID, rangeStr string) ([]byte, error) {
	from, to, err := parseRange(rangeStr)
	if err != nil {
		return nil, err
	}

	sessions, err := s.exportRepo.GetSessionsByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	writer.Write([]string{
		"Session ID", "Purchase Date", "Supplier", "Payment Method",
		"Invoice Reference", "Product", "Category", "Quantity", "Unit Price", "Subtotal",
	})

	// Rows
	for _, session := range sessions {
		for _, item := range session.ProductItems {
			subtotal := float64(item.Quantity) * item.UnitPrice
			writer.Write([]string{
				session.ID.String(),
				session.PurchaseDate.Format("2006-01-02"),
				session.SupplierName,
				string(session.PaymentMethod),
				derefString(session.InvoiceReference),
				item.Name,
				item.Category,
				fmt.Sprintf("%d", item.Quantity),
				fmt.Sprintf("%.2f", item.UnitPrice),
				fmt.Sprintf("%.2f", subtotal),
			})
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *ExportService) ExportPDF(ctx context.Context, userID uuid.UUID, rangeStr string) ([]byte, error) {
	from, to, err := parseRange(rangeStr)
	if err != nil {
		return nil, err
	}

	sessions, err := s.exportRepo.GetSessionsByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "BizKeeper - Purchase Report")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 8, fmt.Sprintf("Period: %s to %s", from.Format("2006-01-02"), to.Format("2006-01-02")))
	pdf.Ln(10)

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(230, 230, 230)
	headers := []string{"Date", "Supplier", "Payment Method", "Invoice", "Product", "Category", "Qty", "Unit Price", "Subtotal"}
	widths := []float64{22, 35, 28, 35, 40, 28, 12, 25, 25}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	for _, session := range sessions {
		for _, item := range session.ProductItems {
			subtotal := float64(item.Quantity) * item.UnitPrice
			row := []string{
				session.PurchaseDate.Format("2006-01-02"),
				session.SupplierName,
				string(session.PaymentMethod),
				derefString(session.InvoiceReference),
				item.Name,
				item.Category,
				fmt.Sprintf("%d", item.Quantity),
				fmt.Sprintf("%.2f", item.UnitPrice),
				fmt.Sprintf("%.2f", subtotal),
			}
			for i, cell := range row {
				pdf.CellFormat(widths[i], 7, cell, "1", 0, "L", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}