CREATE TABLE IF NOT EXISTS scans (
    id INTEGER PRIMARY KEY,
    created_at TEXT NOT NULL,
    host TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS processes (
    id INTEGER PRIMARY KEY,
    scan_id INTEGER NOT NULL,
    pid INTEGER NOT NULL,
    name TEXT NOT NULL,
    executable_path TEXT NOT NULL,
    memory_footprint INTEGER NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(id)
);

CREATE TABLE IF NOT EXISTS connections (
    id INTEGER PRIMARY KEY,
    scan_id INTEGER NOT NULL,
    process_pid INTEGER NOT NULL,
    protocol TEXT NOT NULL,
    local_addr TEXT NOT NULL,
    remote_addr TEXT NOT NULL,
    state TEXT NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(id)
);

CREATE TABLE IF NOT EXISTS findings (
    id INTEGER PRIMARY KEY,
    scan_id INTEGER NOT NULL,
    process_pid INTEGER NOT NULL,
    code TEXT NOT NULL,
    severity TEXT NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(id)
);

CREATE INDEX IF NOT EXISTS idx_processes_scan_id ON processes(scan_id);
CREATE INDEX IF NOT EXISTS idx_connections_scan_id ON connections(scan_id);
CREATE INDEX IF NOT EXISTS idx_connections_state ON connections(state);
CREATE INDEX IF NOT EXISTS idx_findings_scan_id ON findings(scan_id);
