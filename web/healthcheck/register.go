package healthcheck

// HealthCheck health check
type HealthCheck interface {
	Name() string
	Do() (interface{}, error)
}

var handlers []HealthCheck

// Register register
func Register(args ...HealthCheck) {
	handlers = append(handlers, args...)
}
