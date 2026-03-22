package ports

import (
	"context"

	"hostlens/internal/domain"
)

// ScanRepository persists and loads scan history.
type ScanRepository interface {
	SaveScan(
		ctx context.Context,
		scan domain.Scan,
		processes []domain.Process,
		connections []domain.Connection,
		findings []domain.Finding,
	) (int64, error)

	ListScans(ctx context.Context, limit int) ([]domain.Scan, error)

	GetScan(ctx context.Context, scanID int64) (
		domain.Scan,
		[]domain.Process,
		[]domain.Connection,
		[]domain.Finding,
		error,
	)
}
