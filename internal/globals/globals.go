package globals

import (
	"context"
	"flag"
	"log/slog"
	"mcp-forge/api"
	"mcp-forge/internal/config"
	"os"
)

type ApplicationContext struct {
	Context context.Context
	Logger  *slog.Logger
	Config  *api.Configuration
}

func NewApplicationContext() (*ApplicationContext, error) {

	appCtx := &ApplicationContext{
		Context: context.Background(),
		Logger:  slog.New(slog.NewJSONHandler(os.Stderr, nil)),
	}

	// Parse and store the config
	var configFlag = flag.String("config", "config.yaml", "path to the config file")
	flag.Parse()

	configContent, err := config.ReadFile(*configFlag)
	if err != nil {
		return appCtx, err
	}
	appCtx.Config = &configContent

	//
	return appCtx, nil
}
