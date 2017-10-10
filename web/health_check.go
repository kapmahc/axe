package web

// HealthCheck health check
type HealthCheck interface {
	Do() error
}
