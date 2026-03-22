package database

import (
	"context"
	"database/sql"
	"fmt"
	"hostlens/internal/domain"
	"time"
)

func (r *ScanRepository) insertScanTx(tx *sql.Tx, ctx context.Context, scan domain.Scan) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		"INSERT INTO scans(created_at, host) VALUES(?, ?)",
		scan.CreatedAt.UTC().Format(time.RFC3339),
		scan.Host,
	)
	if err != nil {
		return 0, fmt.Errorf("exec scanStmt: %w", err)
	}

	scanID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert scan: %w", err)
	}

	return scanID, nil
}

func (r *ScanRepository) saveProcessesTx(tx *sql.Tx, ctx context.Context, scanID int64, ps []domain.Process) error {
	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO processes(scan_id, pid, name, executable_path, memory_footprint)
		 VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare processes save: %w", err)
	}
	defer stmt.Close()

	for _, p := range ps {
		_, err :=stmt.ExecContext(
			ctx,
			scanID,
			p.PID,
			p.Name,
			p.ExecutablePath,
			p.MemoryFootprint,
		)

		if err != nil {
			return fmt.Errorf("insert process pid=%d: %w", p.PID, err)
		}
	}

	return nil
}

func (r *ScanRepository) saveConnectionsTx(
	tx *sql.Tx,
	ctx context.Context,
	scanID int64,
	cs []domain.Connection,
) error {
	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO connections(scan_id, process_pid, protocol, local_addr, remote_addr, state)
		 VALUES (?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	for _, c := range cs {
		_, err := stmt.ExecContext(
			ctx,
			scanID,
			c.ProcessPID,
			c.Protocol,
			c.LocalAddr,
			c.RemoteAddr,
			string(c.State),
		)
		if err != nil {
			return fmt.Errorf("exec context: %w", err)
		}
	}

	return nil
}

func (r *ScanRepository) saveFindingsTx(
	tx *sql.Tx,
	ctx context.Context,
	scanID int64,
	findings []domain.Finding,
) error {
	findingStmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO findings(scan_id, process_pid, code, severity, description)
		 VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare finding insert: %w", err)
	}
	defer findingStmt.Close()

	for _, f := range findings {
		_, err := findingStmt.ExecContext(
			ctx,
			scanID,
			f.ProcessPID,
			f.Code,
			f.Severity,
			f.Description,
		)
		if err != nil {
			return fmt.Errorf("insert finding pid=%d code=%s: %w", f.ProcessPID, f.Code, err)
		}
	}

	return nil
}


