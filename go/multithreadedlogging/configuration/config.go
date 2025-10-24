package configuration

// Config holds configuration constants for the logging system
type Config struct {
	DefaultPublishTimeoutMs     int
	DefaultConsumptionTimeoutMs int
	DefaultCapacity             int
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultPublishTimeoutMs:     100,
		DefaultConsumptionTimeoutMs: 100,
		DefaultCapacity:             100,
	}
}
