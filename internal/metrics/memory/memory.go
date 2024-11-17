package memory

import (
	"strconv"

	"sysinfo/internal/metrics"

	"github.com/shirou/gopsutil/v4/mem"
)

type Memory struct{}

func (m *Memory) GetMetrics() ([]metrics.Metric, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return []metrics.Metric{
		{Name: "Memory", Type: metrics.TypeTitle, Value: "Memory"},
		{Name: "total", Type: metrics.TypeByte, Value: strconv.FormatUint(v.Total, 10)},
		{Name: "available", Type: metrics.TypeByte, Value: strconv.FormatUint(v.Available, 10)},
		{Name: "used", Type: metrics.TypeByte, Value: strconv.FormatUint(v.Used, 10)},
		{Name: "usedPercent", Type: metrics.TypePer, Value: strconv.FormatFloat(v.UsedPercent, 'f', -1, 64)},
		{Name: "free", Type: metrics.TypeByte, Value: strconv.FormatUint(v.Free, 10)},
	}, nil
}

func NewMemory() *Memory {
	return &Memory{}
}
