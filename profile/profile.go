package profile

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
	"github.com/widrin/knights/logger"
)

type SystemCollector struct {
	mu sync.RWMutex

	cpuUsage    prometheus.Gauge
	memUsage    prometheus.Gauge
	diskUsage   *prometheus.GaugeVec
	netBytes    *prometheus.CounterVec
	requestHist prometheus.Histogram
	lastCPUTime int64
	lastCPUIdle int64
}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{
		cpuUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		}),
		memUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_mem_usage_bytes",
			Help: "Current memory usage in bytes",
		}),
		diskUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_disk_usage_bytes",
				Help: "Disk space usage by mount point",
			},
			[]string{"mount"},
		),
		netBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "system_network_bytes_total",
				Help: "Network traffic statistics",
			},
			[]string{"direction"},
		),
		requestHist: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Request duration distribution",
				Buckets: prometheus.DefBuckets,
			}),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	c.cpuUsage.Describe(ch)
	c.memUsage.Describe(ch)
	c.diskUsage.Describe(ch)
	c.netBytes.Describe(ch)
	c.requestHist.Describe(ch)
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.collectCPUUsage()
	c.collectMemoryUsage()
	c.collectDiskUsage()
	c.collectNetworkStats()

	c.cpuUsage.Collect(ch)
	c.memUsage.Collect(ch)
	c.diskUsage.Collect(ch)
	c.netBytes.Collect(ch)
	c.requestHist.Collect(ch)
}

// 各子系统数据采集方法
func (c *SystemCollector) collectCPUUsage() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 实现基于/proc/stat的CPU时间采集
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		logger.Error("读取CPU统计失败: %v", err)
		return
	}

	var user, nice, system, idle, iowait, irq, softirq int64
	_, err = fmt.Sscanf(strings.Split(string(contents), "\n")[0],
		"cpu  %d %d %d %d %d %d %d",
		&user, &nice, &system, &idle, &iowait, &irq, &softirq)
	if err != nil {
		logger.Error("解析CPU统计失败: %v", err)
		return
	}

	total := user + nice + system + idle + iowait + irq + softirq
	idleTotal := idle + iowait

	if c.lastCPUTime > 0 {
		deltaTotal := total - (c.lastCPUTime)
		deltaIdle := idleTotal - (c.lastCPUIdle)
		usage := (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100
		c.cpuUsage.Set(usage)
	}

	c.lastCPUTime = total
	c.lastCPUIdle = idleTotal
}

func (c *SystemCollector) collectMemoryUsage() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var info syscall.Statfs_t
	err := syscall.Statfs("/", &info)
	if err != nil {
		logger.Error("获取内存信息失败: %v", err)
		return
	}

	total := info.Blocks * uint64(info.Bsize)
	free := info.Bfree * uint64(info.Bsize)
	c.memUsage.Set(float64(total - free))
}

func (c *SystemCollector) collectDiskUsage() {
	c.mu.Lock()
	defer c.mu.Unlock()

	parts, err := disk.Partitions(false)
	if err != nil {
		logger.Error("获取磁盘分区失败: %v", err)
		return
	}

	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			logger.Debug("获取挂载点%s使用量失败: %v", part.Mountpoint, err)
			continue
		}
		c.diskUsage.WithLabelValues(part.Mountpoint).Set(float64(usage.Used))
	}
}

func (c *SystemCollector) collectNetworkStats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ioCounters, err := net.IOCounters(false); err == nil {
		for _, stats := range ioCounters {
			c.netBytes.WithLabelValues("in").Add(float64(stats.BytesRecv))
			c.netBytes.WithLabelValues("out").Add(float64(stats.BytesSent))
		}
	} else {
		logger.Debug("无法获取gnet流量统计")
	}
}
