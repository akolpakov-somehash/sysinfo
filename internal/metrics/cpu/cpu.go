package cpu

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
)

type Cpu struct{}

func (c *Cpu) GetMetrics() (map[string]string, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"name":      "CPU",
		"cpu":       strconv.Itoa(int(info[0].CPU)),
		"vendorId":  info[0].VendorID,
		"family":    info[0].Family,
		"model":     info[0].Model,
		"stepping":  strconv.Itoa(int(info[0].Stepping)),
		"cores":     strconv.Itoa(int(info[0].Cores)),
		"modelName": info[0].ModelName,
		"mhz":       strconv.FormatFloat(info[0].Mhz, 'f', -1, 64),
		"cacheSize": strconv.Itoa(int(info[0].CacheSize)),
		"flags":     fmt.Sprintf("%v", info[0].Flags),
	}, nil
}

func NewCpu() *Cpu {
	return &Cpu{}
}
