package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"sysinfo/pkg/formatter"
)

type MetricProvider interface {
	GetMetrics() (map[string]string, error)
}

type ProvidersMap map[string]MetricProvider

type metricsFilter []string

func (m *metricsFilter) String() string {
	return fmt.Sprint(*m)
}

func (m *metricsFilter) Set(value string) error {
	if len(*m) > 0 {
		return errors.New("metrics filter already set")
	}
	for _, f := range strings.Split(value, ",") {
		*m = append(*m, f)
	}
	return nil
}

type Dummy1 struct{}

func (d *Dummy1) GetMetrics() (map[string]string, error) {
	return map[string]string{
		"name": "cpu",
		"cpu":  "100%",
	}, nil
}

type Dummy2 struct{}

func (d *Dummy2) GetMetrics() (map[string]string, error) {
	return map[string]string{
		"name": "mem",
		"mem":  "100%",
	}, nil
}

func main() {

	metricsFlag := metricsFilter{}
	flag.Var(&metricsFlag, "filter", "comma-separated list of metrics, available: cpu,mem,net,disk,osinf")
	logLevel := flag.String("log-level", "info", "log level (debug, info, warn, error, fatal, panic)")
	format := flag.String("format", "text", "output format (text, json)")
	flag.Parse()
	if len(metricsFlag) == 0 {
		metricsFlag = metricsFilter{"cpu", "mem", "net", "disk", "osinf"}
	}
	lvl, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse log level")
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(lvl)

	formatter := formatter.NewFormatter(*format)

	providers := ProvidersMap{
		"cpu": &Dummy1{},
		"mem": &Dummy2{},
	}

	m := make(chan map[string]string, len(metricsFlag))
	started := 0
	for _, metric := range metricsFlag {
		provider, ok := providers[metric]
		if !ok {
			log.Warn().Str("metric", metric).Msg("Unknown metric")
			continue
		}
		started++
		go func() {
			metrics, err := provider.GetMetrics()
			if err != nil {
				log.Error().Err(err).Str("metric", metric).Msg("Failed to get metrics")
				m <- map[string]string{"error": err.Error()}
				return
			}
			m <- metrics
		}()
	}

	if started == 0 {
		log.Info().Msg("No metrics to collect")
		return
	}

	allMetrics := make([]map[string]string, started)
	for i := 0; i < started; i++ {
		metrics := <-m
		allMetrics[i] = metrics
	}

	fmt.Println("All metrics:")
	fmt.Println(formatter.Format(allMetrics))
}
