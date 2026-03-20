package workflow

import (
	"context"
	"errors"
	"testing"

	"novel-video-workflow/pkg/providers"
)

func TestExecutor_RunChapterWorkflow_HappyPath(t *testing.T) {
	exec := newExecutorWithMocks(t)
	result, err := exec.RunChapterWorkflow(context.Background(), RunRequest{ChapterID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusSucceeded {
		t.Fatalf("expected success, got %s", result.Status)
	}
}

func TestExecutor_PersistsStepTransitions(t *testing.T) {
	store := NewMemoryRunStorage()
	exec := newExecutorWithStore(t, store)

	_, err := exec.RunChapterWorkflow(context.Background(), RunRequest{ChapterID: 2})
	if err != nil {
		t.Fatal(err)
	}

	run, err := store.Load(2)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != StatusSucceeded {
		t.Fatalf("expected succeeded status, got %s", run.Status)
	}
	if len(run.Artifacts) != 4 {
		t.Fatalf("expected 4 artifacts, got %d", len(run.Artifacts))
	}
}

func TestExecutor_RecordsArtifactMetadata(t *testing.T) {
	store := NewMemoryRunStorage()
	exec := newExecutorWithStore(t, store)

	_, err := exec.RunChapterWorkflow(context.Background(), RunRequest{ChapterID: 3})
	if err != nil {
		t.Fatal(err)
	}

	run, err := store.Load(3)
	if err != nil {
		t.Fatal(err)
	}

	ttsArtifact, ok := run.Artifacts[StepTTS]
	if !ok {
		t.Fatal("expected tts artifact")
	}
	if _, ok := ttsArtifact["audio_path"]; !ok {
		t.Fatal("expected audio_path in tts artifact")
	}
}

func TestExecutor_ReturnsCategorizedErrorOnProviderFailure(t *testing.T) {
	exec := NewExecutor(providers.ProviderBundle{
		TTS:      failingTTSProvider{},
		Subtitle: stubSubtitleProvider{},
		Image:    stubImageProvider{},
		Project:  stubProjectProvider{},
	}, NewMemoryRunStorage())

	_, err := exec.RunChapterWorkflow(context.Background(), RunRequest{ChapterID: 4})
	if err == nil {
		t.Fatal("expected error")
	}

	var providerErr providers.ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T", err)
	}
}

func TestExecutor_ResumeFromFailedStep(t *testing.T) {
	store := NewMemoryRunStorage()

	// First run fails at subtitle step
	exec1 := newExecutorWithMockTTSAndFailingSubtitle(t, store)
	_, err := exec1.RunChapterWorkflow(context.Background(), RunRequest{ChapterID: 5})
	if err == nil {
		t.Fatal("expected error from first run")
	}

	// Second run should resume from subtitle step
	exec2 := newExecutorWithMocksAndStore(t, store)
	result, err := exec2.RunChapterWorkflow(context.Background(), RunRequest{ChapterID: 5})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusSucceeded {
		t.Fatalf("expected succeeded after resume, got %s", result.Status)
	}
}

// Helper functions

func newExecutorWithMocks(t *testing.T) *Executor {
	t.Helper()
	return NewExecutor(providers.ProviderBundle{
		TTS:      stubTTSProvider{},
		Subtitle: stubSubtitleProvider{},
		Image:    stubImageProvider{},
		Project:  stubProjectProvider{},
	}, NewMemoryRunStorage())
}

func newExecutorWithStore(t *testing.T, storage Storage) *Executor {
	t.Helper()
	return NewExecutor(providers.ProviderBundle{
		TTS:      stubTTSProvider{},
		Subtitle: stubSubtitleProvider{},
		Image:    stubImageProvider{},
		Project:  stubProjectProvider{},
	}, storage)
}

func newExecutorWithMockTTSAndFailingSubtitle(t *testing.T, storage Storage) *Executor {
	t.Helper()
	return NewExecutor(providers.ProviderBundle{
		TTS:      stubTTSProvider{},
		Subtitle: failingSubtitleProvider{},
		Image:    stubImageProvider{},
		Project:  stubProjectProvider{},
	}, storage)
}

func newExecutorWithMocksAndStore(t *testing.T, storage Storage) *Executor {
	t.Helper()
	return NewExecutor(providers.ProviderBundle{
		TTS:      stubTTSProvider{},
		Subtitle: stubSubtitleProvider{},
		Image:    stubImageProvider{},
		Project:  stubProjectProvider{},
	}, storage)
}

// Stub providers

type stubTTSProvider struct{}

func (stubTTSProvider) Name() string { return "stub-tts" }
func (stubTTSProvider) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "stub-tts", Severity: providers.SeverityInfo, Message: "ready"}
}
func (stubTTSProvider) Generate(req providers.TTSRequest) (providers.TTSResult, error) {
	return providers.TTSResult{AudioPath: "/tmp/audio.wav"}, nil
}

type stubSubtitleProvider struct{}

func (stubSubtitleProvider) Name() string { return "stub-subtitle" }
func (stubSubtitleProvider) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "stub-subtitle", Severity: providers.SeverityInfo, Message: "ready"}
}
func (stubSubtitleProvider) Generate(req providers.SubtitleRequest) (providers.SubtitleResult, error) {
	return providers.SubtitleResult{SubtitlePath: "/tmp/subtitle.srt", Format: "srt"}, nil
}

type stubImageProvider struct{}

func (stubImageProvider) Name() string { return "stub-image" }
func (stubImageProvider) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "stub-image", Severity: providers.SeverityInfo, Message: "ready"}
}
func (stubImageProvider) Generate(req providers.ImageRequest) (providers.ImageResult, error) {
	return providers.ImageResult{ImagePaths: []string{"/tmp/image.png"}}, nil
}

type stubProjectProvider struct{}

func (stubProjectProvider) Name() string { return "stub-project" }
func (stubProjectProvider) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "stub-project", Severity: providers.SeverityInfo, Message: "ready"}
}
func (stubProjectProvider) Generate(req providers.ProjectRequest) (providers.ProjectResult, error) {
	return providers.ProjectResult{ProjectPath: "/tmp/project.json"}, nil
}

type failingTTSProvider struct{}

func (failingTTSProvider) Name() string { return "failing-tts" }
func (failingTTSProvider) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "failing-tts", Severity: providers.SeverityInfo, Message: "ready"}
}
func (failingTTSProvider) Generate(req providers.TTSRequest) (providers.TTSResult, error) {
	return providers.TTSResult{}, providers.NewProviderError(providers.CategoryExecutionError, "tts failed", nil)
}

type failingSubtitleProvider struct{}

func (failingSubtitleProvider) Name() string { return "failing-subtitle" }
func (failingSubtitleProvider) HealthCheck() providers.HealthCheckResult {
	return providers.HealthCheckResult{Provider: "failing-subtitle", Severity: providers.SeverityInfo, Message: "ready"}
}
func (failingSubtitleProvider) Generate(req providers.SubtitleRequest) (providers.SubtitleResult, error) {
	return providers.SubtitleResult{}, providers.NewProviderError(providers.CategoryExecutionError, "subtitle failed", nil)
}