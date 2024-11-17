package cpu

import (
	"fmt"
	"strconv"

	"sysinfo/internal/metrics"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPU is a struct that implements the metrics.Provider interface for CPU metrics.
type CPU struct{}

// GetMetrics returns CPU metrics.
func (c *CPU) GetMetrics() ([]metrics.Metric, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("cpu.Info() failed: %w", err)
	}
	if len(info) == 0 {
		return nil, fmt.Errorf("no CPU information available")
	}

	return []metrics.Metric{
		{Name: "CPU", Type: metrics.TypeTitle, Value: "CPU"},
		{Name: "vendorId", Type: metrics.TypeStr, Value: info[0].VendorID},
		{Name: "cores", Type: metrics.TypeInt, Value: strconv.Itoa(int(info[0].Cores))},
		{Name: "modelName", Type: metrics.TypeStr, Value: info[0].ModelName},
		{Name: "mhz", Type: metrics.TypeStr, Value: fmt.Sprintf("%.2f Mhz", info[0].Mhz)},
		{Name: "cacheSize", Type: metrics.TypeStr, Value: fmt.Sprintf("%d Kb", info[0].CacheSize)},
	}, nil
}

// NewCPU creates a new CPU instance.
func NewCPU() *CPU {
	return &CPU{}
}
