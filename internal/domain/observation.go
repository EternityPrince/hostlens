package domain

import "time"

type Scan struct {
	ID        int64
	CreatedAt time.Time
	Host      string
}

type Process struct {
	ID              int64
	ScanID          int64
	PID             int
	Name            string
	ExecutablePath  string
	MemoryFootprint uint64
}

type ConnectionState string

const (
	StateListen      ConnectionState = "LISTEN"
	StateEstablished ConnectionState = "ESTABLISHED"
	StateOther       ConnectionState = "OTHER"
)

type Connection struct {
	ID         int64
	ScanID     int64
	ProcessPID int
	Protocol   string
	LocalAddr  string
	RemoteAddr string
	State      ConnectionState
}

type Finding struct {
	ID          int64
	ScanID      int64
	ProcessPID  int
	Code        string
	Severity    string
	Description string
}
