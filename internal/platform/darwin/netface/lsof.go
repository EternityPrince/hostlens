package netface

import (
	"bufio"
	"context"
	"fmt"
	"hostlens/internal/domain"
	"hostlens/internal/ports"
	"os/exec"
	"strconv"
	"strings"
)

type Provider struct {}

func NewProvider() *Provider{
	return &Provider{}
}

var _ ports.ConnectionProvider = (*Provider)(nil)

func (p *Provider) ListConnections(ctx context.Context) ([]domain.Connection, error) {
	cmd := exec.CommandContext(ctx, "lsof", "-nP", "-iTCP")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error lsof exec: %w", err)
	}

	connections, err := parseLsofOutput(string(out))
	if err != nil {
		return nil, fmt.Errorf("error parsing output: %w", err)
	}

	return connections, nil
}

func parseLsofOutput(out string) ([]domain.Connection, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	connections := make([]domain.Connection, 0)

	firstLine := true
	for scanner.Scan(){
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if firstLine{
			firstLine = false
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("error parse pid: %w", err)
		}

		nameField := strings.Join(fields[8:], " ")

		localAddr, remoteAddr, state := parseNameField(nameField)

		connections = append(connections, domain.Connection{
			ProcessPID: pid,
			Protocol:   strings.ToUpper(fields[7]),
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			State:      state,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan lsof output: %w", err)
	}


	return connections, nil
}

func parseNameField(nf string) (string, string, domain.ConnectionState) {
	if strings.HasSuffix(nf, " (ESTABLISHED)") {
		base := strings.TrimSuffix(nf, " (ESTABLISHED)")
		parts := strings.SplitN(base, "->", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], domain.StateEstablished
		}
		return base, "", domain.StateEstablished
	}

	if strings.HasSuffix(nf, " (LISTEN)") {
		local := strings.TrimSuffix(nf, " (LISTEN)")
		return local, "", domain.StateListen
	}

	if strings.Contains(nf, "->") {
		parts := strings.SplitN(nf, "->", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], domain.StateOther
		}
	}
	
	return nf, "", domain.StateOther
}
