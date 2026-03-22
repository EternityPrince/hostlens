package database

import (
	"hostlens/internal/domain"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestScanRepository_SaveScanAndGetScan(t *testing.T) {
	ctx, _, repo := newTestRepository(t)
	inputScan := domain.Scan{
		CreatedAt: time.Date(2026, 3, 21, 15, 30, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
		Host:      "test-host",
	}

	inputProcesses := []domain.Process{
		{
			PID:             42,
			Name:            "login",
			ExecutablePath:  "/usr/bin/login",
			MemoryFootprint: 512,
		},
		{
			PID:             1,
			Name:            "launchd",
			ExecutablePath:  "/sbin/launchd",
			MemoryFootprint: 1024,
		},
	}

	inputConnections := []domain.Connection{
		{
			ProcessPID: 42,
			Protocol:   "TCP",
			LocalAddr:  "127.0.0.1:8080",
			RemoteAddr: "127.0.0.1:52344",
			State:      domain.StateEstablished,
		},
		{
			ProcessPID: 1,
			Protocol:   "TCP",
			LocalAddr:  "0.0.0.0:22",
			RemoteAddr: "",
			State:      domain.StateListen,
		},
	}

	inputFindings := []domain.Finding{
		{
			ProcessPID:  42,
			Code:        "EXPOSED_LISTEN",
			Severity:    "info",
			Description: "Process is listening on a non-loopback address",
		},
		{
			ProcessPID:  1,
			Code:        "SSH_EXPOSED",
			Severity:    "warn",
			Description: "SSH is reachable from non-loopback interfaces",
		},
	}

	scanID, err := repo.SaveScan(ctx, inputScan, inputProcesses, inputConnections, inputFindings)
	if err != nil {
		t.Fatalf("SaveScan returned error: %v", err)
	}

	gotScan, gotProcesses, gotConnections, gotFindings, err := repo.GetScan(ctx, scanID)
	if err != nil {
		t.Fatalf("GetScan returned error: %v", err)
	}

	if gotScan.ID != scanID {
		t.Fatalf("expected scan ID %d, got %d", scanID, gotScan.ID)
	}

	wantCreatedAt := inputScan.CreatedAt.UTC()
	if !gotScan.CreatedAt.Equal(wantCreatedAt) {
		t.Fatalf("expected CreatedAt %v, got %v", wantCreatedAt, gotScan.CreatedAt)
	}

	if gotScan.Host != inputScan.Host {
		t.Fatalf("expected Host %q, got %q", inputScan.Host, gotScan.Host)
	}

	wantProcesses := []domain.Process{
		{
			ScanID:          scanID,
			PID:             1,
			Name:            "launchd",
			ExecutablePath:  "/sbin/launchd",
			MemoryFootprint: 1024,
		},
		{
			ScanID:          scanID,
			PID:             42,
			Name:            "login",
			ExecutablePath:  "/usr/bin/login",
			MemoryFootprint: 512,
		},
	}

	if len(gotProcesses) != len(wantProcesses) {
		t.Fatalf("expected %d processes, got %d", len(wantProcesses), len(gotProcesses))
	}

	for i := range gotProcesses {
		if gotProcesses[i].ID == 0 {
			t.Fatalf("expected process %d to have database ID assigned", i)
		}
		gotProcesses[i].ID = 0
	}

	if !reflect.DeepEqual(gotProcesses, wantProcesses) {
		t.Fatalf("processes mismatch:\nwant: %#v\ngot:  %#v", wantProcesses, gotProcesses)
	}

	wantConnections := []domain.Connection{
		{
			ScanID:     scanID,
			ProcessPID: 1,
			Protocol:   "TCP",
			LocalAddr:  "0.0.0.0:22",
			RemoteAddr: "",
			State:      domain.StateListen,
		},
		{
			ScanID:     scanID,
			ProcessPID: 42,
			Protocol:   "TCP",
			LocalAddr:  "127.0.0.1:8080",
			RemoteAddr: "127.0.0.1:52344",
			State:      domain.StateEstablished,
		},
	}

	if len(gotConnections) != len(wantConnections) {
		t.Fatalf("expected %d connections, got %d", len(wantConnections), len(gotConnections))
	}

	for i := range gotConnections {
		if gotConnections[i].ID == 0 {
			t.Fatalf("expected connection %d to have database ID assigned", i)
		}
		gotConnections[i].ID = 0
	}

	if !reflect.DeepEqual(gotConnections, wantConnections) {
		t.Fatalf("connections mismatch:\nwant: %#v\ngot:  %#v", wantConnections, gotConnections)
	}

	wantFindings := []domain.Finding{
		{
			ScanID:      scanID,
			ProcessPID:  1,
			Code:        "SSH_EXPOSED",
			Severity:    "warn",
			Description: "SSH is reachable from non-loopback interfaces",
		},
		{
			ScanID:      scanID,
			ProcessPID:  42,
			Code:        "EXPOSED_LISTEN",
			Severity:    "info",
			Description: "Process is listening on a non-loopback address",
		},
	}

	if len(gotFindings) != len(wantFindings) {
		t.Fatalf("expected %d findings, got %d", len(wantFindings), len(gotFindings))
	}

	for i := range gotFindings {
		if gotFindings[i].ID == 0 {
			t.Fatalf("expected finding %d to have database ID assigned", i)
		}
		gotFindings[i].ID = 0
	}

	if !reflect.DeepEqual(gotFindings, wantFindings) {
		t.Fatalf("findings mismatch:\nwant: %#v\ngot:  %#v", wantFindings, gotFindings)
	}
}

func TestScanRepository_ListScans_OrderAndLimit(t *testing.T) {
	ctx, _, repo := newTestRepository(t)

	scansToSave := []domain.Scan{
		{
			CreatedAt: time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC),
			Host:      "host-1",
		},
		{
			CreatedAt: time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC),
			Host:      "host-2",
		},
		{
			CreatedAt: time.Date(2026, 3, 21, 8, 0, 0, 0, time.UTC),
			Host:      "host-3",
		},
	}

	for _, scan := range scansToSave {
		if _, err := repo.SaveScan(ctx, scan, nil, nil, nil); err != nil {
			t.Fatalf("SaveScan returned error: %v", err)
		}
	}

	got, err := repo.ListScans(ctx, 2)
	if err != nil {
		t.Fatalf("ListScans returned error: %v", err)
	}

	want := []domain.Scan{
		{
			ID:        3,
			CreatedAt: time.Date(2026, 3, 21, 8, 0, 0, 0, time.UTC),
			Host:      "host-3",
		},
		{
			ID:        2,
			CreatedAt: time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC),
			Host:      "host-2",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListScans mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestScanRepository_GetScan_NotFound(t *testing.T) {
	ctx, _, repo := newTestRepository(t)

	_, _, _, _, err := repo.GetScan(ctx, 99999)
	if err == nil {
		t.Fatalf("expected error for missing scan")
	}

	if !strings.Contains(err.Error(), "no such scan") {
		t.Fatalf("expected missing scan error, got %v", err)
	}
}
