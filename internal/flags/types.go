package flags

type Command string

const (
	CommandHelp   Command = "help"
	CommandScan   Command = "scan"
	CommandScans  Command = "scans"
	CommandShow   Command = "show"
	CommandLatest Command = "latest"
)

type ScanOptions struct {
	Save  bool
	JSON  bool
	PID   int
	Name  string
	State string
}

type ScansOptions struct {
	Limit int
	JSON  bool
}

type ShowOptions struct {
	ScanID int64
	JSON   bool
}

type LatestOptions struct {
	JSON bool
}

type Options struct {
	Command Command
	DBPath  string

	Scan   ScanOptions
	Scans  ScansOptions
	Show   ShowOptions
	Latest LatestOptions
}
