package storage

type StreamLoader interface {
	// SetModel(model any)
	// GetStream() chan any
	Start() (<-chan any, <-chan error)
}