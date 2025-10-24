package constant

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	WARN
	ERROR
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
