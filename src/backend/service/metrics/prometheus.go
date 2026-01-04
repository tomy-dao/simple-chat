package metrics

import (
	"local/model"
	"local/service/common"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics holds all Prometheus metrics collectors
type PrometheusMetrics struct {
	params *common.Params

	// Database metrics
	usersTotal         prometheus.Gauge
	conversationsTotal prometheus.Gauge
	messagesTotal      prometheus.Gauge

	// Application metrics
	uptimeSeconds prometheus.Gauge

	// System metrics
	goroutinesCount prometheus.Gauge
	memoryAllocBytes prometheus.Gauge
	memoryTotalBytes prometheus.Gauge
	memorySysBytes   prometheus.Gauge
	gcRunsTotal      prometheus.Gauge
}

// NewPrometheusMetrics creates a new Prometheus metrics collector
func NewPrometheusMetrics(params *common.Params) *PrometheusMetrics {
	pm := &PrometheusMetrics{
		params: params,

		usersTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_users_total",
			Help: "Total number of registered users",
		}),

		conversationsTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_conversations_total",
			Help: "Total number of conversations",
		}),

		messagesTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_messages_total",
			Help: "Total number of messages",
		}),

		uptimeSeconds: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_uptime_seconds",
			Help: "Application uptime in seconds",
		}),

		goroutinesCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_goroutines",
			Help: "Current number of goroutines",
		}),

		memoryAllocBytes: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_memory_alloc_bytes",
			Help: "Bytes of allocated heap objects",
		}),

		memoryTotalBytes: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_memory_total_bytes",
			Help: "Cumulative bytes allocated for heap objects",
		}),

		memorySysBytes: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_memory_sys_bytes",
			Help: "Total bytes of memory obtained from the OS",
		}),

		gcRunsTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "simple_chat_gc_runs_total",
			Help: "Total number of completed GC cycles",
		}),
	}

	// Start background collector
	go pm.collectMetrics()

	return pm
}

// collectMetrics periodically collects and updates all metrics
func (pm *PrometheusMetrics) collectMetrics() {
	startTime := time.Now()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// Collect immediately on start
	pm.updateMetrics(startTime)

	for range ticker.C {
		pm.updateMetrics(startTime)
	}
}

// updateMetrics updates all Prometheus metrics
func (pm *PrometheusMetrics) updateMetrics(startTime time.Time) {
	reqCtx := model.NewRequestContext(nil)

	// Count users using repository method
	if userCount, err := pm.params.Repo.UserRepo.Count(reqCtx); err == nil {
		pm.usersTotal.Set(float64(userCount))
	}

	// Count conversations using repository method
	if conversationCount, err := pm.params.Repo.ConversationRepo.Count(reqCtx); err == nil {
		pm.conversationsTotal.Set(float64(conversationCount))
	}

	// Count messages using repository method
	if messageCount, err := pm.params.Repo.MessageRepo.Count(reqCtx); err == nil {
		pm.messagesTotal.Set(float64(messageCount))
	}

	// Update application uptime
	uptime := time.Since(startTime)
	pm.uptimeSeconds.Set(uptime.Seconds())

	// Update system metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	pm.goroutinesCount.Set(float64(runtime.NumGoroutine()))
	pm.memoryAllocBytes.Set(float64(m.Alloc))
	pm.memoryTotalBytes.Set(float64(m.TotalAlloc))
	pm.memorySysBytes.Set(float64(m.Sys))
	pm.gcRunsTotal.Set(float64(m.NumGC))
}
