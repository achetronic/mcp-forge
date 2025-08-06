package globals

import (
	"context"
	"log/slog"
	"os"
)

var (
	Context = context.Background()
	Logger  = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)
