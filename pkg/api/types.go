package api

import "time"

// ProjectResponse represents a project resource in API responses.
type ProjectResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ChapterResponse represents a chapter resource in API responses.
type ChapterResponse struct {
	ID                 uint      `json:"id"`
	ProjectID          uint      `json:"project_id"`
	Title              string    `json:"title"`
	Content            string    `json:"content,omitempty"`
	SegmentationPrompt string    `json:"segmentation_prompt,omitempty"`
	WorkflowParams     string    `json:"workflow_params,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// SceneResponse represents a scene resource in API responses.
type SceneResponse struct {
	ID          uint      `json:"id"`
	ChapterID   uint      `json:"chapter_id"`
	SceneNumber int       `json:"scene_number"`
	Content     string    `json:"content"`
	ImagePrompt string    `json:"image_prompt,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WorkflowRunResponse represents a workflow run resource in API responses.
type WorkflowRunResponse struct {
	ID             uint      `json:"id"`
	ProjectID      uint      `json:"project_id"`
	ChapterID      uint      `json:"chapter_id"`
	Status         string    `json:"status"`
	CurrentStep    string    `json:"current_step,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	ArtifactMeta   string    `json:"artifact_meta,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	CompletedAt    time.Time `json:"completed_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SystemCheckResultResponse represents a health check result in API responses.
type SystemCheckResultResponse struct {
	Provider string `json:"provider"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// ErrorResponse represents a structured error in API responses.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a generic success response.
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}