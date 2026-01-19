package health

import (
	"context"
	"testing"

	"github.com/tian841224/stock-bot/internal/application/port"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type mockHealthChecker struct {
	checkDatabaseFunc   func(ctx context.Context) port.HealthStatus
	checkAPIFunc        func(ctx context.Context, apiName string) port.HealthStatus
	checkSyncStatusFunc func(ctx context.Context) port.SyncHealthStatus
	checkResourcesFunc  func(ctx context.Context) port.ResourceHealthStatus
}

func (m *mockHealthChecker) CheckDatabase(ctx context.Context) port.HealthStatus {
	if m.checkDatabaseFunc != nil {
		return m.checkDatabaseFunc(ctx)
	}
	return port.HealthStatus{Status: "healthy", Message: "Connected", ResponseTime: 5}
}

func (m *mockHealthChecker) CheckAPI(ctx context.Context, apiName string) port.HealthStatus {
	if m.checkAPIFunc != nil {
		return m.checkAPIFunc(ctx, apiName)
	}
	return port.HealthStatus{Status: "healthy", Message: "API responsive", ResponseTime: 100}
}

func (m *mockHealthChecker) CheckSyncStatus(ctx context.Context) port.SyncHealthStatus {
	if m.checkSyncStatusFunc != nil {
		return m.checkSyncStatusFunc(ctx)
	}
	return port.SyncHealthStatus{Status: "healthy"}
}

func (m *mockHealthChecker) CheckResources(ctx context.Context) port.ResourceHealthStatus {
	if m.checkResourcesFunc != nil {
		return m.checkResourcesFunc(ctx)
	}
	return port.ResourceHealthStatus{Status: "healthy", MemoryUsageMB: 100, Goroutines: 10, CPUCores: 4}
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Error(msg string, fields ...logger.Field) {}
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Debug(msg string, fields ...logger.Field) {}
func (m *mockLogger) Panic(msg string, fields ...logger.Field) {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Field) {}
func (m *mockLogger) Sync() error                              { return nil }

func TestHealthCheckUsecase_GetHealthStatus_AllHealthy(t *testing.T) {
	mockChecker := &mockHealthChecker{}
	usecase := NewHealthCheckUsecase(mockChecker, "test-service", "1.0.0", &mockLogger{})

	response, err := usecase.GetHealthStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status healthy, got %s", response.Status)
	}

	if !response.OverallHealthy {
		t.Error("Expected overall_healthy to be true")
	}

	if response.Service != "test-service" {
		t.Errorf("Expected service test-service, got %s", response.Service)
	}

	if response.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", response.Version)
	}
}

func TestHealthCheckUsecase_GetHealthStatus_DatabaseUnhealthy(t *testing.T) {
	mockChecker := &mockHealthChecker{
		checkDatabaseFunc: func(ctx context.Context) port.HealthStatus {
			return port.HealthStatus{Status: "unhealthy", Message: "Connection failed", ResponseTime: 0}
		},
	}
	usecase := NewHealthCheckUsecase(mockChecker, "test-service", "1.0.0", &mockLogger{})

	response, err := usecase.GetHealthStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Status != "unhealthy" {
		t.Errorf("Expected status unhealthy, got %s", response.Status)
	}

	if response.OverallHealthy {
		t.Error("Expected overall_healthy to be false")
	}
}

func TestHealthCheckUsecase_GetHealthStatus_APIDegraded(t *testing.T) {
	mockChecker := &mockHealthChecker{
		checkAPIFunc: func(ctx context.Context, apiName string) port.HealthStatus {
			return port.HealthStatus{Status: "degraded", Message: "Slow response", ResponseTime: 5000}
		},
	}
	usecase := NewHealthCheckUsecase(mockChecker, "test-service", "1.0.0", &mockLogger{})

	response, err := usecase.GetHealthStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Status != "degraded" {
		t.Errorf("Expected status degraded, got %s", response.Status)
	}

	if !response.OverallHealthy {
		t.Error("Expected overall_healthy to be true (degraded but not critical)")
	}
}

func TestHealthCheckUsecase_GetHealthStatus_ResourcesDegraded(t *testing.T) {
	mockChecker := &mockHealthChecker{
		checkResourcesFunc: func(ctx context.Context) port.ResourceHealthStatus {
			return port.ResourceHealthStatus{Status: "degraded", MemoryUsageMB: 2000, Goroutines: 1500, CPUCores: 4}
		},
	}
	usecase := NewHealthCheckUsecase(mockChecker, "test-service", "1.0.0", &mockLogger{})

	response, err := usecase.GetHealthStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Status != "degraded" {
		t.Errorf("Expected status degraded, got %s", response.Status)
	}

	if !response.OverallHealthy {
		t.Error("Expected overall_healthy to be true (degraded but not critical)")
	}
}

func TestHealthCheckUsecase_DetermineOverallStatus(t *testing.T) {
	usecase := &healthCheckUsecase{}

	tests := []struct {
		name            string
		dbStatus        port.HealthStatus
		finmindStatus   port.HealthStatus
		fugleStatus     port.HealthStatus
		syncStatus      port.SyncHealthStatus
		resourceStatus  port.ResourceHealthStatus
		expectedStatus  string
		expectedHealthy bool
	}{
		{
			name:            "All healthy",
			dbStatus:        port.HealthStatus{Status: "healthy"},
			finmindStatus:   port.HealthStatus{Status: "healthy"},
			fugleStatus:     port.HealthStatus{Status: "healthy"},
			syncStatus:      port.SyncHealthStatus{Status: "healthy"},
			resourceStatus:  port.ResourceHealthStatus{Status: "healthy"},
			expectedStatus:  "healthy",
			expectedHealthy: true,
		},
		{
			name:            "Database unhealthy",
			dbStatus:        port.HealthStatus{Status: "unhealthy"},
			finmindStatus:   port.HealthStatus{Status: "healthy"},
			fugleStatus:     port.HealthStatus{Status: "healthy"},
			syncStatus:      port.SyncHealthStatus{Status: "healthy"},
			resourceStatus:  port.ResourceHealthStatus{Status: "healthy"},
			expectedStatus:  "unhealthy",
			expectedHealthy: false,
		},
		{
			name:            "API degraded",
			dbStatus:        port.HealthStatus{Status: "healthy"},
			finmindStatus:   port.HealthStatus{Status: "degraded"},
			fugleStatus:     port.HealthStatus{Status: "healthy"},
			syncStatus:      port.SyncHealthStatus{Status: "healthy"},
			resourceStatus:  port.ResourceHealthStatus{Status: "healthy"},
			expectedStatus:  "degraded",
			expectedHealthy: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, healthy := usecase.determineOverallStatus(
				tt.dbStatus,
				tt.finmindStatus,
				tt.fugleStatus,
				tt.syncStatus,
				tt.resourceStatus,
			)

			if status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, status)
			}

			if healthy != tt.expectedHealthy {
				t.Errorf("Expected healthy %v, got %v", tt.expectedHealthy, healthy)
			}
		})
	}
}
