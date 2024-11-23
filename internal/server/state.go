package server

import "context"

// State represents the current state of a server.
type State int

const (
	StateStarting State = iota
	StateRunning
	StateShuttingDown
	StateStopped
)

// Server interface defines the common operations for all servers.
type Server interface {
	Start() error
	Stop(ctx context.Context) error
	IsReady() bool
	HealthCheck() bool
	State() State
}
