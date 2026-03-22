package flags

import "fmt"

type ArgBuilder struct{}

func NewArgBuilder() *ArgBuilder {
	return &ArgBuilder{}
}

func (b *ArgBuilder) Build(args []string) (Options, error) {
	dbPath, err := defaultDBPath()
	if err != nil {
		return Options{}, err
	}

	raw, err := parseArgs(args, dbPath)
	if err != nil {
		return Options{}, err
	}

	return buildOptions(raw)
}

func buildOptions(raw rawOptions) (Options, error) {
	opts := Options{DBPath: raw.DBPath}

	switch raw.Command {
	case "help":
		opts.Command = CommandHelp
		return opts, nil

	case "scan":
		opts.Command = CommandScan
		opts.Scan = ScanOptions{
			Save:  raw.ScanSave,
			JSON:  raw.ScanJSON,
			PID:   raw.ScanPID,
			Name:  raw.ScanName,
			State: raw.ScanState,
		}

		return opts, nil

	case "scans":
		if raw.ScansLimit <= 0 {
			return Options{}, fmt.Errorf("limit must be > 0")
		}
		opts.Command = CommandScans
		opts.Scans = ScansOptions{
			Limit: raw.ScansLimit,
			JSON:  raw.ScansJSON,
		}
		return opts, nil

	case "show":
		if raw.ShowScanID <= 0 {
			return Options{}, UsageError("scan ID must be > 0")
		}
		opts.Command = CommandShow
		opts.Show = ShowOptions{
			ScanID: raw.ShowScanID,
			JSON:   raw.ShowJSON,
		}
		return opts, nil

	case "latest":
		opts.Command = CommandLatest
		opts.Latest = LatestOptions{
			JSON: raw.LatestJSON,
		}
		return opts, nil

	default:
		return Options{}, fmt.Errorf("unsupported command: %s", raw.Command)
	}
}
