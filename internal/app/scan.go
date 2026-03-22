package app

import (
	"context"
	"fmt"
	"hostlens/internal/analyze"
	"hostlens/internal/domain"
	"os"
	"time"
)

func (a *App) runScan(ctx context.Context) error {
	processes, err := a.processSource.ListProcesses(ctx)
	if err != nil {
		return fmt.Errorf("error process: %w", err)
	}

	connections, err := a.connectionSource.ListConnections(ctx)
	if err != nil {
		return fmt.Errorf("error get list connections: %w", err)
	}

	findings := analyze.BuildFindings(connections)

	fmt.Printf("found %d processes\n", len(processes))
	fmt.Printf("found %d tcp connections\n", len(connections))
	fmt.Printf("found %d findings\n", len(findings))

	limit := 10
	limit = min(limit, len(processes))
	_ = limit

	//	for _, p := range processes[:limit]{
	//		fmt.Printf("pid=%d name=%q path=%q\n", p.PID, p.Name, p.ExecutablePath)
	//	}

	if !a.opts.Scan.Save {
		fmt.Println("scan result was not saved; rerun with --save to store it in the database")
		return nil
	}

	host, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error get hostname: %w", err)
	}

	scanID, err := a.scanRepo.SaveScan(
		ctx,
		domain.Scan{
			CreatedAt: time.Now(),
			Host:      host,
		},
		processes,
		connections,
		findings,
	)
	if err != nil {
		return fmt.Errorf("bad commit on save scan: %w", err)
	}

	fmt.Printf("saved scan id=%d\n", scanID)

	return nil
}
