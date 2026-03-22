package app

import (
	"context"
	"fmt"
)

func (a *App) runScans (ctx context.Context) (error) {
	scans, err := a.scanRepo.ListScans(ctx, a.opts.Scans.Limit)
	if err != nil {
		return fmt.Errorf("get scans error: %w", err)
	}

	if len(scans) <= 0 {
		fmt.Println("no saved scans")
		return nil
	}

	for _, s := range scans {
		fmt.Printf("id=%d created_at=%s host=%q\n",
			s.ID,
			s.CreatedAt.Format("2006-01-02 15:04:05"),
			s.Host,
		)
	}
	return nil
}
