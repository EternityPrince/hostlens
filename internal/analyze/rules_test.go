package analyze

import (
	"hostlens/internal/domain"
	"reflect"
	"testing"
)

func TestBuildFindings(t *testing.T) {
	input := []domain.Connection{
		{
			ProcessPID: 101,
			LocalAddr:  "0.0.0.0:443",
			State:      domain.StateListen,
		},
		{
			ProcessPID: 202,
			LocalAddr:  "127.0.0.1:8080",
			State:      domain.StateListen,
		},
		{
			ProcessPID: 303,
			LocalAddr:  "localhost:9000",
			State:      domain.StateListen,
		},
		{
			ProcessPID: 404,
			LocalAddr:  "[::1]:7000",
			State:      domain.StateListen,
		},
		{
			ProcessPID: 505,
			LocalAddr:  "192.168.1.10:5432",
			State:      domain.StateEstablished,
		},
	}

	got := BuildFindings(input)
	want := []domain.Finding{
		{
			ProcessPID:  101,
			Code:        FindingCodeExposedListen,
			Severity:    SeverityInfo,
			Description: "Process is listening on a non-loopback address",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildFindings mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestIsLoopbackListen(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want bool
	}{
		{name: "ipv4 loopback", addr: "127.0.0.1:80", want: true},
		{name: "localhost hostname", addr: "localhost:3000", want: true},
		{name: "ipv6 bracketed loopback", addr: "[::1]:9000", want: true},
		{name: "ipv6 compact loopback", addr: "::1:9000", want: true},
		{name: "wildcard listen", addr: "0.0.0.0:22", want: false},
		{name: "private interface", addr: "192.168.1.4:22", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isLoopbackListen(tc.addr); got != tc.want {
				t.Fatalf("isLoopbackListen(%q) = %v, want %v", tc.addr, got, tc.want)
			}
		})
	}
}
