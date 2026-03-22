package analyze

import (
	"hostlens/internal/domain"
	"strings"
)

const (
	FindingCodeExposedListen = "EXPOSED_LISTEN"
	SeverityInfo             = "info"
)

func BuildFindings(cs []domain.Connection) []domain.Finding {
	findings := make([]domain.Finding, 0)

	for _, c := range cs {
		if c.State != domain.StateListen {
			continue
		}

		if isLoopbackListen(c.LocalAddr) {
			continue
		}

		findings = append(findings, domain.Finding{
			ProcessPID:  c.ProcessPID,
			Code:        FindingCodeExposedListen,
			Severity:    SeverityInfo,
			Description: "Process is listening on a non-loopback address",
		})
	}

	return findings
}

func isLoopbackListen(addr string) bool {
	return (
		strings.HasPrefix(addr, "127.") ||
		strings.HasPrefix(addr, "localhost:") ||
		strings.HasPrefix(addr, "[::1]:") ||
		strings.HasPrefix(addr, "::1:"))
}
