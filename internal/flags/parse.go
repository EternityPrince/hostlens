package flags

import (
	stdflag "flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type rawOptions struct {
	Command string
	DBPath  string

	ScanSave  bool
	ScanJSON  bool
	ScanPID   int
	ScanName  string
	ScanState string

	ScansLimit int
	ScansJSON  bool

	ShowScanID int64
	ShowJSON   bool

	LatestJSON bool
}

// parseArgs reads CLI flags into a raw structure.
// It does not apply business validation.
func parseArgs(args []string, defaultDB string) (rawOptions, error) {
	var raw rawOptions

	if len(args) == 0 {
		raw.Command = string(CommandHelp)
		raw.DBPath = defaultDB
		return raw, nil
	}
	raw.Command = args[0]

	switch raw.Command {
	case "help", "-h", "--help":
		raw.Command = string(CommandHelp)
		raw.DBPath = defaultDB
		return raw, nil
	}

	switch raw.Command {
	case "scan":
		fs := stdflag.NewFlagSet("scan", stdflag.ContinueOnError)
		fs.SetOutput(os.Stderr)

		fs.StringVar(&raw.DBPath, "db-path", defaultDB, "path to SQLite database")
		fs.BoolVar(&raw.ScanSave, "save", false, "save scan result to database")
		fs.BoolVar(&raw.ScanJSON, "json", false, "render scan as JSON")
		fs.IntVar(&raw.ScanPID, "pid", 0, "filter by process PID")
		fs.StringVar(&raw.ScanName, "name", "", "filter by process name")
		fs.StringVar(&raw.ScanState, "state", "", "filter by connection state")

		if err := fs.Parse(args[1:]); err != nil {
			return rawOptions{}, err
		}

	case "scans":
		fs := stdflag.NewFlagSet("scans", stdflag.ContinueOnError)
		fs.SetOutput(os.Stderr)

		fs.StringVar(&raw.DBPath, "db-path", defaultDB, "path to SQLite database")
		fs.IntVar(&raw.ScansLimit, "limit", 10, "number of saved scans to show")
		fs.BoolVar(&raw.ScansJSON, "json", false, "render saved scans as JSON")

		if err := fs.Parse(args[1:]); err != nil {
			return rawOptions{}, err
		}

	case "latest":
		fs := stdflag.NewFlagSet("latest", stdflag.ContinueOnError)
		fs.SetOutput(os.Stderr)

		fs.StringVar(&raw.DBPath, "db-path", defaultDB, "path to SQLite database")
		fs.BoolVar(&raw.LatestJSON, "json", false, "render latest saved scan as JSON")

		if err := fs.Parse(args[1:]); err != nil {
			return rawOptions{}, err
		}

	case "show":
		fs := stdflag.NewFlagSet("show", stdflag.ContinueOnError)
		fs.SetOutput(os.Stderr)

		fs.StringVar(&raw.DBPath, "db-path", defaultDB, "path to SQLite database")
		fs.Int64Var(&raw.ShowScanID, "scan-id", 0, "saved scan ID")
		fs.BoolVar(&raw.ShowJSON, "json", false, "render saved scan as JSON")

		showArgs := args[1:]
		if len(showArgs) > 0 && !strings.HasPrefix(showArgs[0], "-") {
			showArgs = append([]string{"--scan-id", showArgs[0]}, showArgs[1:]...)
		}

		if err := fs.Parse(showArgs); err != nil {
			return rawOptions{}, err
		}

		if fs.NArg() > 1 {
			return rawOptions{}, UsageError("show accepts at most one positional scan ID")
		}

		if raw.ShowScanID == 0 && fs.NArg() == 1 {
			scanID, err := strconv.ParseInt(fs.Arg(0), 10, 64)
			if err != nil {
				return rawOptions{}, UsageError(fmt.Sprintf("invalid scan ID %q", fs.Arg(0)))
			}
			raw.ShowScanID = scanID
		}

	case "help":
		raw.DBPath = defaultDB
		return raw, nil

	default:
		return rawOptions{}, UsageError(fmt.Sprintf("unknown command %q", raw.Command))
	}
	return raw, nil
}
