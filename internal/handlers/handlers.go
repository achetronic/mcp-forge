package handlers

import "mcp-forge/internal/globals"

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
