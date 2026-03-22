package ports

import (
	"context"

	"hostlens/internal/domain"
)

// ProcessProvider loads process metadata from the operating system.
type ProcessProvider interface {
	ListProcesses(ctx context.Context) ([]domain.Process, error)
}
