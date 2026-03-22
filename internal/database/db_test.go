package database

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestOpenAndInit(t *testing.T) {
	dbPath := newTestDBPath(t)

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.Init(ctx); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	rawDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}
	defer rawDB.Close()

	for _, object := range []struct {
		kind string
		name string
	}{
		{kind: "table", name: "scans"},
		{kind: "table", name: "processes"},
		{kind: "table", name: "connections"},
		{kind: "table", name: "findings"},
		{kind: "index", name: "idx_processes_scan_id"},
		{kind: "index", name: "idx_connections_scan_id"},
		{kind: "index", name: "idx_connections_state"},
		{kind: "index", name: "idx_findings_scan_id"},
	} {
		var name string
		err := rawDB.QueryRowContext(ctx, `
			SELECT name
			FROM sqlite_master
			WHERE type = ? AND name = ?
		`, object.kind, object.name).Scan(&name)
		if err != nil {
			t.Fatalf("%s %q not found: %v", object.kind, object.name, err)
		}
	}
}

func TestSplitSQLStatement(t *testing.T) {
	src := `
		CREATE TABLE scans(id INTEGER PRIMARY KEY);

		CREATE INDEX idx_scans_id ON scans(id);
	`

	got := splitSQLStatement(src)
	want := []string{
		"CREATE TABLE scans(id INTEGER PRIMARY KEY);",
		"CREATE INDEX idx_scans_id ON scans(id);",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitSQLStatement mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}
