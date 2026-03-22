package app

import (
	"context"
	"fmt"
)

func (a *App) runLatest(ctx context.Context) (error) {
	scans, err := a.scanRepo.ListScans(ctx, 1)
	if err != nil {
		return fmt.Errorf("error scans: %w", err)
	}

	if len(scans) < 1 {
		fmt.Println("no saved scans")
		return nil
	}

	scan, processes, connections, findings, err := a.scanRepo.GetScan(ctx, scans[0].ID)
	if err != nil {
		return fmt.Errorf("get scan problem: %w", err)
	}
	ss := Snapshot{s: scan, p: processes, c: connections, f: findings}

	a.showSnapshot(ss, false)

	return nil
}
