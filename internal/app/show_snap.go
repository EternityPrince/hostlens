package app

import (
	"fmt"
	"hostlens/internal/domain"
)

type Snapshot struct {
	s domain.Scan
	p []domain.Process
	c []domain.Connection
	f []domain.Finding
}

func (a *App) showSnapshot(s Snapshot, details bool) {
	fmt.Printf(
		"id=%d created_at=%s host=%q processes=%d connections=%d findings=%d\n",
		s.s.ID,
		s.s.CreatedAt.Format("2006-01-02 15:04:05"),
		s.s.Host,
		len(s.p),
		len(s.c),
		len(s.f),
	)
	if !details {
		return
	}

	pLimit := min(10, len(s.p))
	for _, p := range s.p[:pLimit] {
		fmt.Printf("pid=%d name=%q path=%q\n", p.PID, p.Name, p.ExecutablePath)
	}

	cLimit := min(10, len(s.c))
	for _, c := range s.c[:cLimit] {
		fmt.Printf(
			"conn pid=%d proto=%s local=%q remote=%q state=%q\n",
			c.ProcessPID,
			c.Protocol,
			c.LocalAddr,
			c.RemoteAddr,
			c.State,
		)
	}

	fLimit := min(10, len(s.f))
	for _, f := range s.f[:fLimit] {
		fmt.Printf(
			"finding pid=%d code=%q severity=%q desc=%q\n",
			f.ProcessPID,
			f.Code,
			f.Severity,
			f.Description,
		)
	}
}
