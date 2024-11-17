package metrics

// Metric category identifiers.
const (
	CPU    = "cpu"
	Mem    = "mem"
	Net    = "net"
	Disk   = "disk"
	OSInfo = "osinf"
)

// Metric data types.
type MetricType int

// The list of metric types.
const (
	TypeInt MetricType = iota
	TypeByte
	TypeStr
	TypePer
	TypeAny
	TypeTitle
)

// Metric represents a single metric.
type Metric struct {
	Name  string
	Type  MetricType
	Value string
}

// MetricProvider defines an interface for fetching metrics.
type MetricProvider interface {
	// GetMetrics retrieves a slice of metrics or an error if the operation fails.
	GetMetrics() ([]Metric, error)
}

// ProvidersMap is a map of metric providers.
type ProvidersMap map[string]MetricProvider
