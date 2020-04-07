// Copyright (c) 2018 Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"resenje.org/logging"

	"github.com/prometheus/client_golang/prometheus"
)

// Handler increments Prometheus counter metrics partitioned by log level.
type Handler struct {
	Counter *Counter
	Level   logging.Level
}

// NewHandler returns new Handler instance.
func NewHandler(counter *Counter, level logging.Level) (h *Handler) {
	return &Handler{
		Counter: counter,
		Level:   level,
	}
}

// Handle increments Prometheus counter metrics based on
// record log level.
func (h *Handler) Handle(record *logging.Record) error {
	switch record.Level {
	case logging.EMERGENCY:
		h.Counter.emergencyCounter.Inc()
	case logging.ALERT:
		h.Counter.alertCounter.Inc()
	case logging.CRITICAL:
		h.Counter.criticalCounter.Inc()
	case logging.ERROR:
		h.Counter.errorCounter.Inc()
	case logging.WARNING:
		h.Counter.warningCounter.Inc()
	case logging.NOTICE:
		h.Counter.noticeCounter.Inc()
	case logging.INFO:
		h.Counter.infoCounter.Inc()
	case logging.DEBUG:
		h.Counter.debugCounter.Inc()
	}
	return nil
}

// HandleError does nothing for this handler.
func (h *Handler) HandleError(err error) error {
	return nil
}

// GetLevel returns current level for this handler.
func (h *Handler) GetLevel() logging.Level {
	return h.Level
}

// Close does nothing for this handler.
func (h *Handler) Close() error {
	return nil
}

// Counter holds Prometheus counters.
type Counter struct {
	vector           *prometheus.CounterVec
	emergencyCounter prometheus.Counter
	alertCounter     prometheus.Counter
	criticalCounter  prometheus.Counter
	errorCounter     prometheus.Counter
	warningCounter   prometheus.Counter
	noticeCounter    prometheus.Counter
	infoCounter      prometheus.Counter
	debugCounter     prometheus.Counter
}

// CounterOptions holds options for NewCounter constructor.
type CounterOptions struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

// NewCounter creates new Counter instance.
// Options value can be nil.
func NewCounter(options *CounterOptions) (c *Counter) {
	if options == nil {
		options = new(CounterOptions)
	}
	if options.Subsystem == "" {
		options.Subsystem = "logging"
	}
	if options.Name == "" {
		options.Name = "messages_total"
	}
	if options.Help == "" {
		options.Help = "Number of log messages processed, partitioned by log level."
	}
	vector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: options.Namespace,
			Subsystem: options.Subsystem,
			Name:      options.Name,
			Help:      options.Help,
		},
		[]string{"level"},
	)
	return &Counter{
		vector:           vector,
		emergencyCounter: vector.WithLabelValues("emergency"),
		alertCounter:     vector.WithLabelValues("alert"),
		criticalCounter:  vector.WithLabelValues("critical"),
		errorCounter:     vector.WithLabelValues("error"),
		warningCounter:   vector.WithLabelValues("warning"),
		noticeCounter:    vector.WithLabelValues("notice"),
		infoCounter:      vector.WithLabelValues("info"),
		debugCounter:     vector.WithLabelValues("debug"),
	}
}

// Metrics retuns all Prometheus metrics that
// should be registered.
func (c *Counter) Metrics() (cs []prometheus.Collector) {
	return []prometheus.Collector{c.vector}
}
