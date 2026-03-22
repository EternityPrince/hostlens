package flags

import "fmt"

func UsageError(msg string) error {
	return fmt.Errorf("%s\n\n%s", msg, usageText())
}

func usageText() string {
	return `Usage:
  hostlens help
  hostlens scan [--db-path PATH] [--save] [--json] [--pid N] [--name NAME] [--state STATE]
  hostlens scans [--db-path PATH] [--limit N] [--json]
  hostlens show [--db-path PATH] (--scan-id N | N) [--json]
  hostlens latest [--db-path PATH] [--json]
`
}
