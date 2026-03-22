package flags

import (
	"fmt"
	"os"
	"path/filepath"
)

const appName = "hostlens"

// defaultDBPath returns the default SQLite database path.
// It creates an app-specific runtime directory in the user's config location.

func defaultDBPath() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	appDir := filepath.Join(baseDir, appName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", fmt.Errorf("error when making directories: %w", err)
	}

	return filepath.Join(appDir, "hostlens.db"), nil
}
