package proc

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

// Provider loads process data from macOS.
// This first version uses the "ps" command as a simple live backend.
type Provider struct{}

// NewProvider creates a new process provider.
func NewProvider() *Provider {
	return &Provider{}
}

// Compile-time interface check.
// If Provider stops matching ports.ProcessProvider, the build will fail.
var _ ports.ProcessProvider = (*Provider)(nil)

// ListProcesses returns a live list of processes.
// This first implementation fills PID and Name only.
// ExecutablePath and MemoryFootprint will be added later.
func (p *Provider) ListProcesses(ctx context.Context) ([]domain.Process, error) {
	cmd := exec.CommandContext(ctx, "ps", "-axo", "pid=,comm=")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}
	
	process, err := parsePSOutput(string(out))
	if err != nil {
		return nil, fmt.Errorf("error parsing PSO: %w", err)
	}

	return process, nil
}

func parsePSOutput(src string) ([]domain.Process, error) {
	scanner := bufio.NewScanner(strings.NewReader(src))
	process := make([]domain.Process, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		firstSpace := strings.IndexAny(line, " \t")
		if firstSpace == -1 {
			continue
		}

		pidPart := strings.TrimSpace(line[:firstSpace])
		commandPath := strings.TrimSpace(line[firstSpace+1:])
		if commandPath == "" {
			continue
		}

		pid, err := strconv.Atoi(pidPart)
		if err != nil {
			return nil, fmt.Errorf("parse PID error: %w", err)
		}


		process = append(process, domain.Process{
			PID: pid,
			Name: shortName(commandPath),
			ExecutablePath: commandPath,
		})

		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error scanner: %w", err)
		}

	return process, nil
}

	func shortName(path string) string {
		if path == "" {
			return ""
		}

		idx := strings.LastIndex(path, "/")
		if idx == -1 {
			return path
		}
		return path[idx+1:]
}


