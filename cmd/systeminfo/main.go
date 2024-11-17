package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"sysinfo/internal/metrics"
	"sysinfo/internal/metrics/cpu"
	"sysinfo/internal/metrics/disk"
	"sysinfo/internal/metrics/memory"

	"sysinfo/pkg/formatter"
)

// metricsFilter is a custom flag type for metrics filter.
type metricsFilter []string

// String returns a string representation of the metrics filter.
func (m *metricsFilter) String() string {
	return fmt.Sprint(*m)
}

// Set sets the metrics filter value.
func (m *metricsFilter) Set(value string) error {
	for _, f := range strings.Split(value, ",") {
		*m = append(*m, f)
	}
	return nil
}

// defaultMetrics is a list of default metrics to collect.
var defaultMetrics = metricsFilter{metrics.CPU, metrics.Mem, metrics.Net, metrics.Disk, metrics.OSInfo}

// main is the entry point of the program.
func main() {
	metricsFlag, logLevel, format := parseFlags()
	initializeLogging(logLevel)
	formatter, err := initializeFormatter(format)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize formatter")
		return
	}
	providers := initializeProviders()
	allMetrics, err := collectMetrics(providers, metricsFlag)
	if err != nil {
		log.Error().Err(err).Msg("failed to collect metrics")
		return
	}
	outputMetrics(allMetrics, formatter)
}

// parseFlags parses command line flags.
func parseFlags() (metricsFilter, string, string) {
	metricsFlag := metricsFilter{}
	flag.Var(&metricsFlag, "filter", "comma-separated list of metrics, available: cpu,mem,net,disk,osinf")
	logLevel := flag.String("log-level", "info", "log level (debug, info, warn, error, fatal, panic)")
	format := flag.String("format", formatter.TextFormat, "output format (text, json)")
	flag.Parse()
	if len(metricsFlag) == 0 {
		metricsFlag = append(metricsFilter{}, defaultMetrics...)
	}
	return metricsFlag, *logLevel, *format
}

// initializeLogging initializes logging.
func initializeLogging(logLevel string) {
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(lvl)
}

// initializeProviders initializes metric providers.
func initializeProviders() metrics.ProvidersMap {
	return metrics.ProvidersMap{
		metrics.CPU:  cpu.NewCPU(),
		metrics.Mem:  memory.NewMemory(),
		metrics.Disk: disk.NewDisk(),
	}
}

// collectMetrics collects metrics from providers.
func collectMetrics(providers metrics.ProvidersMap, metricsFlag metricsFilter) ([]metrics.MetricGroup, error) {
	m := make(chan metrics.MetricGroup, len(metricsFlag))
	started := 0
	for _, metric := range metricsFlag {
		metric := metric
		provider, ok := providers[metric]
		if !ok {
			log.Warn().Str("metric", metric).Msg("unknown metric")
			continue
		}
		started++
		go func(metric string, provider metrics.MetricProvider) {
			metricsResult, err := provider.GetMetrics()
			if err != nil {
				log.Error().Err(err).Str("metric", metric).Msg("failed to get metrics")
				m <- metrics.MetricGroup{Metrics: []metrics.Metric{{Name: metric, Type: metrics.TypeStr, Value: "failed to get metrics"}}}
				return
			}
			m <- metricsResult
		}(metric, provider)
	}

	if started == 0 {
		return nil, fmt.Errorf("no metrics to collect")
	}

	allMetrics := make([]metrics.MetricGroup, started)
	for i := 0; i < started; i++ {
		metricsResult := <-m
		allMetrics[i] = metricsResult
	}
	return allMetrics, nil
}

// initializeFormatter initializes formatter.
func initializeFormatter(format string) (formatter.Formatter, error) {
	return formatter.NewFormatter(format)
}

// outputMetrics outputs metrics.
func outputMetrics(metrics []metrics.MetricGroup, formatter formatter.Formatter) {
	formattedMetrics, err := formatter.Format(metrics, "")
	if err != nil {
		log.Error().Err(err).Msg("failed to format metrics")
		return
	}
	fmt.Println("All metrics:\n----------------")
	fmt.Println(formattedMetrics)
}
