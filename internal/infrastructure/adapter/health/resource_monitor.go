package health

import (
	"context"
	"runtime"

	"github.com/tian841224/stock-bot/internal/application/port"
)

type resourceMonitor struct{}

func NewResourceMonitor() *resourceMonitor {
	return &resourceMonitor{}
}

func (m *resourceMonitor) CheckResources(ctx context.Context) port.ResourceHealthStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	memoryMB := int64(mem.Alloc / 1024 / 1024)
	goroutines := runtime.NumGoroutine()
	cpuCores := runtime.NumCPU()

	status := "healthy"
	if memoryMB > 1024 {
		status = "degraded"
	}
	if goroutines > 1000 {
		status = "degraded"
	}

	return port.ResourceHealthStatus{
		Status:        status,
		MemoryUsageMB: memoryMB,
		Goroutines:    goroutines,
		CPUCores:      cpuCores,
	}
}
