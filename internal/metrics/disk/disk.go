package disk

import (
	"fmt"
	"strconv"

	"sysinfo/internal/metrics"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/disk"
)

// Disk is a struct that implements the metrics.Provider interface for disk metrics.
type Disk struct{}

// GetMetrics returns disk metrics.
func (d *Disk) GetMetrics() (metrics.MetricGroup, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return metrics.MetricGroup{}, fmt.Errorf("disk.Partitions() failed: %w", err)
	}

	diskMetrics := metrics.MetricGroup{
		Title:  "Disks",
		Groups: []metrics.MetricGroup{},
	}
	for _, part := range parts {
		if part.Fstype == "" {
			continue
		}
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			diskMetrics.Groups = append(diskMetrics.Groups, metrics.MetricGroup{
				Title: part.Mountpoint,
				Metrics: []metrics.Metric{
					{Name: "error", Type: metrics.TypeStr, Value: fmt.Sprintf("failed to get disk usage: %v", err)},
				},
			})
			log.Error().Err(err).Str("mountpoint", part.Mountpoint).Msg("failed to get disk usage")
			continue
		}

		partitionMetrics := []metrics.Metric{
			{Name: "device", Type: metrics.TypeStr, Value: part.Device},
			{Name: "mountpoint", Type: metrics.TypeStr, Value: part.Mountpoint},
			{Name: "fstype", Type: metrics.TypeStr, Value: part.Fstype},
			{Name: "total", Type: metrics.TypeByte, Value: strconv.FormatUint(usage.Total, 10)},
			{Name: "free", Type: metrics.TypeByte, Value: strconv.FormatUint(usage.Free, 10)},
			{Name: "used", Type: metrics.TypeByte, Value: strconv.FormatUint(usage.Used, 10)},
			{Name: "usedPercent", Type: metrics.TypePer, Value: strconv.FormatFloat(usage.UsedPercent, 'f', -1, 64)},
		}
		diskMetrics.Groups = append(diskMetrics.Groups, metrics.MetricGroup{
			Title:   part.Mountpoint,
			Metrics: partitionMetrics,
		})

	}

	return diskMetrics, nil
}

// NewDisk creates a new Disk instance.
func NewDisk() *Disk {
	return &Disk{}
}
