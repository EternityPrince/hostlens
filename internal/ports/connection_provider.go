package ports

import (
	"context"

	"hostlens/internal/domain"
)

// ConnectionProvider loads network connections from the operating system.
type ConnectionProvider interface {
	ListConnections(ctx context.Context) ([]domain.Connection, error)
}
