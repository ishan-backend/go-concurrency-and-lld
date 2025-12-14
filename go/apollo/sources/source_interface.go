package sources

// DataSource defines interface for reading data, adapter in implementation
type DataSource interface {
	Read() ([]map[string]interface{}, error)
	Close() error
}
