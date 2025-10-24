package constant

// SupportedLogMedium represents the different mediums where logs can be written
type SupportedLogMedium int

const (
	FILE SupportedLogMedium = iota
	NETWORK
	DB
)

// String returns the string representation of SupportedLogMedium
func (m SupportedLogMedium) String() string {
	switch m {
	case FILE:
		return "FILE"
	case NETWORK:
		return "NETWORK"
	case DB:
		return "DB"
	default:
		return "UNKNOWN"
	}
}
