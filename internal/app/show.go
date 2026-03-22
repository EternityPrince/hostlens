package app

import (
	"context"
	"fmt"
)

func (a *App) runShow(ctx context.Context) error {
	scan, processes, connections, findings, err := a.scanRepo.GetScan(ctx, a.opts.Show.ScanID)
	ss := Snapshot{s: scan, p: processes, c: connections, f: findings}
	if err != nil {
		return fmt.Errorf("error get scan by ID: %w", err)
	}
	a.showSnapshot(ss, true)
	return nil
}
