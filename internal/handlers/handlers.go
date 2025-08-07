package handlers

import "mcp-server-template/internal/globals"

type HandlersManagerDependencies struct {
	AppCtx *globals.ApplicationContext
}

type HandlersManager struct {
	dependencies HandlersManagerDependencies
}

func NewHandlersManager(deps HandlersManagerDependencies) *HandlersManager {
	return &HandlersManager{
		dependencies: deps,
	}
}
