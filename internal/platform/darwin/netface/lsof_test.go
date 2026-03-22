package netface

import (
	"reflect"
	"testing"

	"hostlens/internal/domain"
)

func TestParseNameField(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantLocal  string
		wantRemote string
		wantState  domain.ConnectionState
	}{
		{
			name:      "listen socket",
			input:     "127.0.0.1:8080 (LISTEN)",
			wantLocal: "127.0.0.1:8080",
			wantState: domain.StateListen,
		},
		{
			name:       "established connection",
			input:      "127.0.0.1:8080->127.0.0.1:52344 (ESTABLISHED)",
			wantLocal:  "127.0.0.1:8080",
			wantRemote: "127.0.0.1:52344",
			wantState:  domain.StateEstablished,
		},
		{
			name:       "other connection",
			input:      "10.0.0.5:12345->10.0.0.10:443",
			wantLocal:  "10.0.0.5:12345",
			wantRemote: "10.0.0.10:443",
			wantState:  domain.StateOther,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			local, remote, state := parseNameField(tc.input)
			if local != tc.wantLocal {
				t.Fatalf("expected local addr %q, got %q", tc.wantLocal, local)
			}
			if remote != tc.wantRemote {
				t.Fatalf("expected remote addr %q, got %q", tc.wantRemote, remote)
			}
			if state != tc.wantState {
				t.Fatalf("expected state %q, got %q", tc.wantState, state)
			}
		})
	}
}

func TestParseLsofOutput(t *testing.T) {
	input := `COMMAND   PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
nginx     100 root    6u  IPv4 0x01                 0t0  TCP 127.0.0.1:8080 (LISTEN)
nginx     100 root    7u  IPv4 0x02                 0t0  TCP 127.0.0.1:8080->127.0.0.1:52344 (ESTABLISHED)
bad       nope root   8u  IPv4 0x03                 0t0  TCP 10.0.0.1:443 (LISTEN)
short     300 root
`

	_, err := parseLsofOutput(input)
	if err == nil {
		t.Fatalf("expected parseLsofOutput to fail on invalid PID")
	}
}

func TestParseLsofOutputSkipsShortRows(t *testing.T) {
	input := `COMMAND   PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
short     300 root
nginx     100 root    6u  IPv4 0x01                 0t0  TCP 127.0.0.1:8080 (LISTEN)
`

	got, err := parseLsofOutput(input)
	if err != nil {
		t.Fatalf("parseLsofOutput returned error: %v", err)
	}

	want := []domain.Connection{
		{
			ProcessPID: 100,
			Protocol:   "TCP",
			LocalAddr:  "127.0.0.1:8080",
			RemoteAddr: "",
			State:      domain.StateListen,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseLsofOutput mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}
