package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"novel-video-workflow/pkg/providers"
	"novel-video-workflow/pkg/workflow"

	"github.com/gin-gonic/gin"
)

func TestSystemCheckAPI_ReturnsProviderHealthChecks(t *testing.T) {
	router := gin.New()
	bundle := providers.ProviderBundle{
		TTS:      mockTTSProviderForSystemCheck{},
		Subtitle: mockSubtitleProviderForSystemCheck{},
		Image:    mockImageProviderForSystemCheck{},
		Project:  mockProjectProviderForSystemCheck{},
	}
	api := NewSystemCheckAPI(bundle)
	api.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/system/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Fatalf("expected status success, got %s", resp.Status)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	results, ok := data["results"].([]interface{})
	if !ok {
		t.Fatal("expected results array")
	}

	if len(results) != 4 {
		t.Fatalf("expected 4 provider results, got %d", len(results))
	}
}

func TestWorkflowRunAPI_StartsNewRun(t *testing.T) {
	router := gin.New()
	storage := workflow.NewMemoryRunStorage()
	exec := workflow.NewExecutor(providers.ProviderBundle{
		TTS:      mockTTSProviderForRun{},
		Subtitle: mockSubtitleProviderForRun{},
		Image:    mockImageProviderForRun{},
		Project:  mockProjectProviderForRun{},
	}, storage)
	api := NewWorkflowRunAPI(exec, storage)
	api.RegisterRoutes(router)

	body := `{"chapter_id": 1}`
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/runs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Status != "success" {
		t.Fatalf("expected status success, got %s", resp.Status)
	}
}

func TestWorkflowRunAPI_GetsRunByID(t *testing.T) {
	router := gin.New()
	storage := workflow.NewMemoryRunStorage()
	exec := workflow.NewExecutor(providers.ProviderBundle{
		TTS:      mockTTSProviderForRun{},
		Subtitle: mockSubtitleProviderForRun{},
		Image:    mockImageProviderForRun{},
		Project:  mockProjectProviderForRun{},
	}, storage)
	api := NewWorkflowRunAPI(exec, storage)
	api.RegisterRoutes(router)

	// Create a run first
	_, _ = exec.RunChapterWorkflow(context.Background(), workflow.RunRequest{ChapterID: 2})

	req := httptest.NewRequest(http.MethodGet, "/api/workflow/runs/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp WorkflowRunResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.ChapterID != 2 {
		t.Fatalf("expected chapter_id 2, got %d", resp.ChapterID)
	}
}

// Mock providers for system check tests
type mockTTSProviderForSystemCheck struct{}

func (mockTTSProviderForSystemCheck) Name() string { return "mock-tts" }
func (mockTTSProviderForSystemCheck) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-tts", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockTTSProviderForSystemCheck) Generate(req providers.TTSRequest) (providers.TTSResult, error) {
	return providers.TTSResult{AudioPath: "/tmp/audio.wav"}, nil
}

type mockSubtitleProviderForSystemCheck struct{}

func (mockSubtitleProviderForSystemCheck) Name() string { return "mock-subtitle" }
func (mockSubtitleProviderForSystemCheck) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-subtitle", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockSubtitleProviderForSystemCheck) Generate(req providers.SubtitleRequest) (providers.SubtitleResult, error) {
	return providers.SubtitleResult{SubtitlePath: "/tmp/subtitle.srt"}, nil
}

type mockImageProviderForSystemCheck struct{}

func (mockImageProviderForSystemCheck) Name() string { return "mock-image" }
func (mockImageProviderForSystemCheck) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-image", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockImageProviderForSystemCheck) Generate(req providers.ImageRequest) (providers.ImageResult, error) {
	return providers.ImageResult{ImagePaths: []string{"/tmp/image.png"}}, nil
}

type mockProjectProviderForSystemCheck struct{}

func (mockProjectProviderForSystemCheck) Name() string { return "mock-project" }
func (mockProjectProviderForSystemCheck) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-project", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockProjectProviderForSystemCheck) Generate(req providers.ProjectRequest) (providers.ProjectResult, error) {
	return providers.ProjectResult{ProjectPath: "/tmp/project.json"}, nil
}

// Mock providers for run tests (same as above, just renamed for clarity)
type mockTTSProviderForRun struct{}

func (mockTTSProviderForRun) Name() string { return "mock-tts" }
func (mockTTSProviderForRun) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-tts", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockTTSProviderForRun) Generate(req providers.TTSRequest) (providers.TTSResult, error) {
	return providers.TTSResult{AudioPath: "/tmp/audio.wav"}, nil
}

type mockSubtitleProviderForRun struct{}

func (mockSubtitleProviderForRun) Name() string { return "mock-subtitle" }
func (mockSubtitleProviderForRun) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-subtitle", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockSubtitleProviderForRun) Generate(req providers.SubtitleRequest) (providers.SubtitleResult, error) {
	return providers.SubtitleResult{SubtitlePath: "/tmp/subtitle.srt"}, nil
}

type mockImageProviderForRun struct{}

func (mockImageProviderForRun) Name() string { return "mock-image" }
func (mockImageProviderForRun) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-image", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockImageProviderForRun) Generate(req providers.ImageRequest) (providers.ImageResult, error) {
	return providers.ImageResult{ImagePaths: []string{"/tmp/image.png"}}, nil
}

type mockProjectProviderForRun struct{}

func (mockProjectProviderForRun) Name() string { return "mock-project" }
func (mockProjectProviderForRun) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "mock-project", Severity: providers.SeverityInfo, Message: "ready"}
}
func (mockProjectProviderForRun) Generate(req providers.ProjectRequest) (providers.ProjectResult, error) {
	return providers.ProjectResult{ProjectPath: "/tmp/project.json"}, nil
}