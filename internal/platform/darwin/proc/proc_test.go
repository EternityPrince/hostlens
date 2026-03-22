package proc

import (
	"reflect"
	"testing"

	"hostlens/internal/domain"
)

func TestParsePSOutput(t *testing.T) {
	input := `
1 /sbin/launchd
42 /usr/bin/login
777 /Applications/Google Chrome.app/Contents/MacOS/Google Chrome
`

	processes, err := parsePSOutput(input)
	if err != nil {
		t.Fatalf("parsePSOutput returned error: %v", err)
	}

	if len(processes) != 3 {
		t.Fatalf("expected 3 processes, got %d", len(processes))
	}

	if processes[0].PID != 1 {
		t.Fatalf("expected PID 1, got %d", processes[0].PID)
	}

	if processes[0].Name != "launchd" {
		t.Fatalf("expected name launchd, got %q", processes[0].Name)
	}

	if processes[2].Name != "Google Chrome" {
		t.Fatalf("expected name Google Chrome, got %q", processes[2].Name)
	}
}

func TestParsePSOutputSkipsMalformedLines(t *testing.T) {
	input := `
bad-line
123
42 /usr/bin/login
`

	got, err := parsePSOutput(input)
	if err != nil {
		t.Fatalf("parsePSOutput returned error: %v", err)
	}

	want := []domain.Process{
		{
			PID:            42,
			Name:           "login",
			ExecutablePath: "/usr/bin/login",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parsePSOutput mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestParsePSOutputInvalidPID(t *testing.T) {
	input := `
abc /usr/bin/login
`

	_, err := parsePSOutput(input)
	if err == nil {
		t.Fatalf("expected parsePSOutput to fail on invalid PID")
	}
}

func TestShortName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "empty", path: "", want: ""},
		{name: "no slash", path: "launchd", want: "launchd"},
		{name: "path", path: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", want: "Google Chrome"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := shortName(tc.path); got != tc.want {
				t.Fatalf("shortName(%q) = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}
