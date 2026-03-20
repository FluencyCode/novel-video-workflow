package api

import "time"

// WebSocket event names (frozen).
const (
	EventWorkflowStarted     = "workflow.started"
	EventWorkflowStepChanged = "workflow.step_changed"
	EventWorkflowLog         = "workflow.log"
	EventWorkflowCompleted   = "workflow.completed"
	EventWorkflowFailed      = "workflow.failed"
	EventSystemCheckUpdated  = "system.check.updated"
)

// WebSocketEventEnvelope represents a WebSocket event envelope.
type WebSocketEventEnvelope struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// NewWebSocketEventEnvelope creates a new WebSocket event envelope.
func NewWebSocketEventEnvelope(eventType string, payload interface{}) WebSocketEventEnvelope {
	return WebSocketEventEnvelope{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// WorkflowStartedPayload represents payload for workflow.started event.
type WorkflowStartedPayload struct {
	ProjectID uint `json:"project_id"`
	ChapterID uint `json:"chapter_id"`
}

// WorkflowStepChangedPayload represents payload for workflow.step_changed event.
type WorkflowStepChangedPayload struct {
	ProjectID uint   `json:"project_id"`
	ChapterID uint   `json:"chapter_id"`
	Step      string `json:"step"`
	Status    string `json:"status"`
}

// WorkflowLogPayload represents payload for workflow.log event.
type WorkflowLogPayload struct {
	ProjectID uint   `json:"project_id"`
	ChapterID uint   `json:"chapter_id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// WorkflowCompletedPayload represents payload for workflow.completed event.
type WorkflowCompletedPayload struct {
	ProjectID    uint   `json:"project_id"`
	ChapterID    uint   `json:"chapter_id"`
	ProjectPath  string `json:"project_path,omitempty"`
	DurationSec  float64 `json:"duration_sec"`
}

// WorkflowFailedPayload represents payload for workflow.failed event.
type WorkflowFailedPayload struct {
	ProjectID    uint   `json:"project_id"`
	ChapterID    uint   `json:"chapter_id"`
	FailedStep   string `json:"failed_step"`
	ErrorMessage string `json:"error_message"`
}

// SystemCheckUpdatedPayload represents payload for system.check.updated event.
type SystemCheckUpdatedPayload struct {
	Results []SystemCheckResultResponse `json:"results"`
	CanStart bool                       `json:"can_start"`
}