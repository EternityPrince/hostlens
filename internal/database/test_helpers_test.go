package database

import (
	"context"
	"path/filepath"
	"testing"
)

func newTestRepository(t *testing.T) (context.Context, *DB, *ScanRepository) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	})

	ctx := context.Background()
	if err := db.Init(ctx); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	return ctx, db, NewScanRepository(db)
}

func newTestDBPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "test.db")
}
