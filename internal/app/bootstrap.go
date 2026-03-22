package app

import (
	"context"
	"fmt"
	"hostlens/internal/database"
	"hostlens/internal/flags"
	"hostlens/internal/platform/darwin/netface"
	"hostlens/internal/platform/darwin/proc"
	"hostlens/internal/ports"
)

type App struct {
	opts             flags.Options
	db               *database.DB
	processSource    ports.ProcessProvider
	scanRepo         ports.ScanRepository
	connectionSource ports.ConnectionProvider
}

// New creates the application, parses CLI arguments,
// opens the database, and initializes the schema.
func New(ctx context.Context, args []string) (*App, error) {
	builder := flags.NewArgBuilder()

	opts, err := builder.Build(args)
	if err != nil {
		return nil, fmt.Errorf("build CLI options: %w", err)
	}

	db, err := database.Open(opts.DBPath)
	if err != nil {
		return nil, fmt.Errorf("db setup error: %w", err)
	}

	if err := db.Init(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db init error: %w", err)
	}

	return &App{
		opts:             opts,
		db:               db,
		processSource:    proc.NewProvider(),
		scanRepo:         database.NewScanRepository(db),
		connectionSource: netface.NewProvider(),
	}, nil
}

// Close releases application resources.
func (a *App) Close() error {
	if a == nil || a.db == nil {
		return nil
	}
	return a.db.Close()
}

// Run executes the selected command.
// For now it only prints the parsed command and DB path.
func (a *App) Run(ctx context.Context) error {
	switch a.opts.Command {
	case flags.CommandHelp:
		fmt.Print(flags.UsageError("").Error())
		return nil
	case flags.CommandScan:
		return a.runScan(ctx)
	case flags.CommandScans:
		return a.runScans(ctx)
	case flags.CommandShow:
		return a.runShow(ctx)
	case flags.CommandLatest:
		return a.runLatest(ctx)

	default:
		return fmt.Errorf("command %q is not implemented yet", a.opts.Command)
	}
}
