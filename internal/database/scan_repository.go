package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hostlens/internal/domain"
	"time"
)

type ScanRepository struct {
	db *DB
}

func NewScanRepository(db *DB) *ScanRepository {
	return &ScanRepository{
		db: db,
	}
}

func (r *ScanRepository) SaveScan(
	ctx context.Context,
	scan domain.Scan,
	processes []domain.Process,
	connections []domain.Connection,
	findings []domain.Finding,
) (int64, error) {

	tx, err := r.db.sql.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	scanID, err := r.insertScanTx(tx, ctx, scan)
	if err != nil {
		return 0, fmt.Errorf("insert scan: %w", err)
	}

	if err := r.saveProcessesTx(tx, ctx, scanID, processes); err != nil {
		return 0, err
	}

	if err := r.saveConnectionsTx(tx, ctx, scanID, connections); err != nil {
		return 0, err
	}

	if err := r.saveFindingsTx(tx, ctx, scanID, findings); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit scan transaction: %w", err)
	}

	return scanID, nil
	}



func (r *ScanRepository) ListScans(ctx context.Context, limit int) ([]domain.Scan, error) {
	rows, err := r.db.sql.QueryContext(
		ctx,
		`SELECT id, created_at, host
		FROM scans
		ORDER BY id DESC
		LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query scans: %w", err)
	}
	defer rows.Close()

	scans := make([]domain.Scan, 0)

	for rows.Next() {
		var scan domain.Scan
		var createdAt string

		if err := rows.Scan(
			&scan.ID, &createdAt, &scan.Host,
		); err != nil {
			return nil, fmt.Errorf("scan scan row: %w", err)
		}

		parsedTime, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}
		scan.CreatedAt = parsedTime
		scans = append(scans, scan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scan rows: %w", err)
	}
	return scans, nil
}

func (r *ScanRepository) listProcessesByScanID(ctx context.Context, scanID int64) ([]domain.Process, error) {
	rows, err := r.db.sql.QueryContext(
		ctx,
		`SELECT id, scan_id, pid, name, executable_path, memory_footprint
		FROM processes
		WHERE scan_id = ?
		ORDER BY pid`,
		scanID,
	)
	if err != nil {
		return nil, fmt.Errorf("query processes by scan id: %w", err)
	}
	defer rows.Close()

	processes := make([]domain.Process, 0)
	for rows.Next() {
		var p domain.Process
		if err := rows.Scan(
			&p.ID,
			&p.ScanID,
			&p.PID,
			&p.Name,
			&p.ExecutablePath,
			&p.MemoryFootprint,
		); err != nil {
			return nil, fmt.Errorf("scan process row: %w", err)
		}
		processes = append(processes, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan process row: %w", err)
	}

	return processes, nil
}

func (r *ScanRepository) listConnectionsByScanID(ctx context.Context, scanID int64) ([]domain.Connection, error) {
	rows, err := r.db.sql.QueryContext(
		ctx,
		`
		SELECT id, scan_id, process_pid, protocol, local_addr, remote_addr, state
		 FROM connections
		 WHERE scan_id = ?
		 ORDER BY process_pid, id
		`,
		scanID,
	)
	if err != nil {
		return nil, fmt.Errorf("error query connections by scanID: %w", err) 
	}
	defer rows.Close()

	connections := make([]domain.Connection, 0)

	for rows.Next() {
		var c domain.Connection
		var state string

		if err := rows.Scan(
			&c.ID,
			&c.ScanID,
			&c.ProcessPID,
			&c.Protocol,
			&c.LocalAddr,
			&c.RemoteAddr,
			&state,
		); err != nil {
			return nil, fmt.Errorf("scan connection row: %w", err)
		}

		c.State = domain.ConnectionState(state)
		connections = append(connections, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan/exec: %w", err)
	}

	return connections, nil 
}

func (r *ScanRepository) GetScan(ctx context.Context, scanID int64) (
	domain.Scan,
	[]domain.Process,
	[]domain.Connection,
	[]domain.Finding,
	error,
) {
	var scan domain.Scan
	var createdAt string

	err := r.db.sql.QueryRowContext(
		ctx,
		"SELECT id, created_at, host FROM scans WHERE id = ?",
		scanID,
	).Scan(&scan.ID, &createdAt, &scan.Host)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Scan{}, nil, nil, nil, fmt.Errorf("no such scan in DB: %w", err)
		}
		return domain.Scan{}, nil, nil, nil, fmt.Errorf("error get scan by ID: %w", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return domain.Scan{}, nil, nil, nil, fmt.Errorf("error parsing time from: %w", err)
	}
	processes, err := r.listProcessesByScanID(ctx, scanID)
	if err != nil {
		return domain.Scan{}, nil, nil, nil, fmt.Errorf("error get processes by scanID: %w", err)
	}

	scan.CreatedAt = parsedTime

	connections, err := r.listConnectionsByScanID(ctx, scanID)
	if err != nil {
		return domain.Scan{}, nil, nil, nil, fmt.Errorf("error get connections by scanID: %w", err)
	}

	fs, err := r.listFindingsByScanID(ctx, scanID)
	if err != nil {
		return domain.Scan{}, nil, nil, nil, fmt.Errorf("error get findings by scanID: %w", err)

	}

	return scan, processes, connections, fs, nil
}

func (r *ScanRepository) listFindingsByScanID(ctx context.Context, scanID int64) ([]domain.Finding, error) {
	rows, err := r.db.sql.QueryContext(
		ctx, 
		`SELECT id, scan_id, process_pid, code, severity, description
		 FROM findings
		 WHERE scan_id = ?
		 ORDER BY process_pid, id`,
		 scanID,
	)
	if err != nil {
		return nil, fmt.Errorf("error query findings: %w", err)
	}
	defer rows.Close()

	fs := make([]domain.Finding, 0)

	for rows.Next() {
		var f domain.Finding
		if err := rows.Scan(
			&f.ID,
			&f.ScanID,
			&f.ProcessPID,
			&f.Code,
			&f.Severity,
			&f.Description,
		); err != nil {
			return nil, fmt.Errorf("scan find: %w", err)
		}

		fs = append(fs, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate finding rows: %w", err)
	}

	return fs, nil
}
