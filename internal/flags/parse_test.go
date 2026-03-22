package flags

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    rawOptions
		wantErr bool
	}{
		{
			name: "scan command",
			args: []string{"scan", "--save", "--json", "--pid", "123", "--name", "chrome", "--state", "LISTEN"},
			want: rawOptions{
				Command:   "scan",
				DBPath:    "/tmp/test.db",
				ScanSave:  true,
				ScanJSON:  true,
				ScanPID:   123,
				ScanName:  "chrome",
				ScanState: "LISTEN",
			},
		},
		{
			name: "scans command uses explicit limit",
			args: []string{"scans", "--limit", "5", "--json"},
			want: rawOptions{
				Command:    "scans",
				DBPath:     "/tmp/test.db",
				ScansLimit: 5,
				ScansJSON:  true,
			},
		},
		{
			name: "latest command uses default db path",
			args: []string{"latest", "--json"},
			want: rawOptions{
				Command:    "latest",
				DBPath:     "/tmp/test.db",
				LatestJSON: true,
			},
		},
		{
			name: "show command",
			args: []string{"show", "--scan-id", "7", "--json"},
			want: rawOptions{
				Command:    "show",
				DBPath:     "/tmp/test.db",
				ShowScanID: 7,
				ShowJSON:   true,
			},
		},
		{
			name: "show command with positional scan id",
			args: []string{"show", "7", "--json"},
			want: rawOptions{
				Command:    "show",
				DBPath:     "/tmp/test.db",
				ShowScanID: 7,
				ShowJSON:   true,
			},
		},
		{
			name: "help command",
			args: []string{"--help"},
			want: rawOptions{
				Command: "help",
				DBPath:  "/tmp/test.db",
			},
		},
		{
			name:    "unknown command",
			args:    []string{"boom"},
			wantErr: true,
		},
		{
			name: "missing command becomes help",
			args: nil,
			want: rawOptions{
				Command: "help",
				DBPath:  "/tmp/test.db",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseArgs(tc.args, "/tmp/test.db")
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected parseArgs to return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseArgs returned error: %v", err)
			}

			if got != tc.want {
				t.Fatalf("parseArgs mismatch:\nwant: %#v\ngot:  %#v", tc.want, got)
			}
		})
	}
}
