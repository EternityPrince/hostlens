package flags

import "testing"

func TestBuildOptions(t *testing.T) {
	tests := []struct {
		name    string
		raw     rawOptions
		want    Options
		wantErr bool
	}{
		{
			name: "help command",
			raw: rawOptions{
				Command: "help",
				DBPath:  "/tmp/test.db",
			},
			want: Options{
				Command: CommandHelp,
				DBPath:  "/tmp/test.db",
			},
		},
		{
			name: "scan command",
			raw: rawOptions{
				Command:   "scan",
				DBPath:    "/tmp/test.db",
				ScanSave:  true,
				ScanJSON:  true,
				ScanPID:   123,
				ScanName:  "chrome",
				ScanState: "ESTABLISHED",
			},
			want: Options{
				Command: CommandScan,
				DBPath:  "/tmp/test.db",
				Scan: ScanOptions{
					Save:  true,
					JSON:  true,
					PID:   123,
					Name:  "chrome",
					State: "ESTABLISHED",
				},
			},
		},
		{
			name: "scans command",
			raw: rawOptions{
				Command:    "scans",
				DBPath:     "/tmp/test.db",
				ScansLimit: 5,
				ScansJSON:  true,
			},
			want: Options{
				Command: CommandScans,
				DBPath:  "/tmp/test.db",
				Scans: ScansOptions{
					Limit: 5,
					JSON:  true,
				},
			},
		},
		{
			name: "show command",
			raw: rawOptions{
				Command:    "show",
				DBPath:     "/tmp/test.db",
				ShowScanID: 7,
				ShowJSON:   true,
			},
			want: Options{
				Command: CommandShow,
				DBPath:  "/tmp/test.db",
				Show: ShowOptions{
					ScanID: 7,
					JSON:   true,
				},
			},
		},
		{
			name: "latest command",
			raw: rawOptions{
				Command:    "latest",
				DBPath:     "/tmp/test.db",
				LatestJSON: true,
			},
			want: Options{
				Command: CommandLatest,
				DBPath:  "/tmp/test.db",
				Latest: LatestOptions{
					JSON: true,
				},
			},
		},
		{
			name: "scans invalid limit",
			raw: rawOptions{
				Command:    "scans",
				DBPath:     "/tmp/test.db",
				ScansLimit: 0,
			},
			wantErr: true,
		},
		{
			name: "show invalid scan id",
			raw: rawOptions{
				Command:    "show",
				DBPath:     "/tmp/test.db",
				ShowScanID: 0,
			},
			wantErr: true,
		},
		{
			name: "unsupported command",
			raw: rawOptions{
				Command: "boom",
				DBPath:  "/tmp/test.db",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildOptions(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected buildOptions to return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("buildOptions returned error: %v", err)
			}

			if got != tc.want {
				t.Fatalf("buildOptions mismatch:\nwant: %#v\ngot:  %#v", tc.want, got)
			}
		})
	}
}
