package server

import "context"

// ServerState represents the current state of a server
type ServerState int

const (
	StateStarting ServerState = iota
	StateRunning
	StateShuttingDown
	StateStopped
)

// Server interface defines the common operations for all servers
type Server interface {
	Start() error
	Stop(ctx context.Context) error
	IsReady() bool
	HealthCheck() bool
	State() ServerState
}
