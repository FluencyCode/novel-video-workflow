package workflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"novel-video-workflow/pkg/providers"
)

// RunRequest initiates a workflow run.
type RunRequest struct {
	ChapterID uint
}

// RunResult captures the outcome of a workflow run.
type RunResult struct {
	ChapterID uint
	Status    Status
}

// Executor runs chapter workflows.
type Executor struct {
	providers providers.ProviderBundle
	storage   Storage
}

// NewExecutor creates a workflow executor.
func NewExecutor(providers providers.ProviderBundle, storage Storage) *Executor {
	return &Executor{
		providers: providers,
		storage:   storage,
	}
}

// RunChapterWorkflow executes the full chapter workflow.
func (e *Executor) RunChapterWorkflow(ctx context.Context, req RunRequest) (RunResult, error) {
	// Load or initialize run state
	run, err := e.storage.LoadByChapterID(req.ChapterID)
	if err != nil {
		// New run
		now := time.Now()
		run = WorkflowRun{
			ChapterID:   req.ChapterID,
			CurrentStep: StepTTS,
			Status:      StatusRunning,
			Artifacts:   map[Step]ArtifactMetadata{},
			StartedAt:   &now,
		}
	} else if run.Status == StatusFailed {
		// Resume from failure
		run = PrepareResume(run)
		run.Status = StatusRunning
	}

	// Save initial state
	if err := e.storage.Save(run); err != nil {
		return RunResult{}, fmt.Errorf("save initial state: %w", err)
	}

	// Execute steps in order
	for _, step := range orderedSteps() {
		// Skip completed steps (from resume)
		if _, ok := run.Artifacts[step]; ok && run.CurrentStep != step {
			continue
		}

		run.CurrentStep = step
		if err := e.storage.Save(run); err != nil {
			return RunResult{}, fmt.Errorf("save step transition: %w", err)
		}

		// Execute step
		artifact, err := e.executeStep(ctx, step, run)
		if err != nil {
			now := time.Now()
			run.Status = StatusFailed
			run.ErrorCategory = categorizeError(err)
			run.ErrorMessage = err.Error()
			run.FinishedAt = &now
			_ = e.storage.Save(run)
			return RunResult{}, err
		}

		run.Artifacts[step] = artifact
		if err := e.storage.Save(run); err != nil {
			return RunResult{}, fmt.Errorf("save artifact: %w", err)
		}
	}

	// Mark succeeded
	now := time.Now()
	run.Status = StatusSucceeded
	run.FinishedAt = &now
	if err := e.storage.Save(run); err != nil {
		return RunResult{}, fmt.Errorf("save final state: %w", err)
	}

	return RunResult{
		ChapterID: run.ChapterID,
		Status:    run.Status,
	}, nil
}

func (e *Executor) executeStep(ctx context.Context, step Step, run WorkflowRun) (ArtifactMetadata, error) {
	switch step {
	case StepTTS:
		return e.executeTTS(ctx, run)
	case StepSubtitle:
		return e.executeSubtitle(ctx, run)
	case StepImage:
		return e.executeImage(ctx, run)
	case StepProject:
		return e.executeProject(ctx, run)
	default:
		return nil, fmt.Errorf("unknown step: %s", step)
	}
}

func (e *Executor) executeTTS(ctx context.Context, run WorkflowRun) (ArtifactMetadata, error) {
	req := providers.TTSRequest{
		// TODO: Load chapter data from database
		Text:      "test text",
		ProjectID: "1",
	}
	result, err := e.providers.TTS.Generate(req)
	if err != nil {
		return nil, err
	}
	return ArtifactMetadata{
		"audio_path": result.AudioPath,
		"duration":   result.Duration,
	}, nil
}

func (e *Executor) executeSubtitle(ctx context.Context, run WorkflowRun) (ArtifactMetadata, error) {
	ttsArtifact := run.Artifacts[StepTTS]
	audioPath, _ := ttsArtifact["audio_path"].(string)

	req := providers.SubtitleRequest{
		AudioPath: audioPath,
		Text:      "test text",
	}
	result, err := e.providers.Subtitle.Generate(req)
	if err != nil {
		return nil, err
	}
	return ArtifactMetadata{
		"subtitle_path": result.SubtitlePath,
		"format":        result.Format,
	}, nil
}

func (e *Executor) executeImage(ctx context.Context, run WorkflowRun) (ArtifactMetadata, error) {
	req := providers.ImageRequest{
		Prompt: "test prompt",
	}
	result, err := e.providers.Image.Generate(req)
	if err != nil {
		return nil, err
	}
	return ArtifactMetadata{
		"image_paths": result.ImagePaths,
	}, nil
}

func (e *Executor) executeProject(ctx context.Context, run WorkflowRun) (ArtifactMetadata, error) {
	req := providers.ProjectRequest{
		ChapterDir: "/tmp/chapter",
	}
	result, err := e.providers.Project.Generate(req)
	if err != nil {
		return nil, err
	}
	return ArtifactMetadata{
		"project_path": result.ProjectPath,
	}, nil
}

func categorizeError(err error) string {
	var providerErr providers.ProviderError
	if errors.As(err, &providerErr) {
		return string(providerErr.Category)
	}
	return "unknown"
}