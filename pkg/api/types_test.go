package api

import (
	"reflect"
	"testing"
)

func TestProjectResponse_HasFrozenCoreFields(t *testing.T) {
	resp := ProjectResponse{}
	typ := reflect.TypeOf(resp)

	fields := []string{"ID", "Title", "CreatedAt", "UpdatedAt"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in ProjectResponse", field)
		}
	}
}

func TestChapterResponse_HasFrozenCoreFields(t *testing.T) {
	resp := ChapterResponse{}
	typ := reflect.TypeOf(resp)

	fields := []string{"ID", "ProjectID", "Title", "Content", "CreatedAt", "UpdatedAt"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in ChapterResponse", field)
		}
	}
}

func TestSceneResponse_HasFrozenCoreFields(t *testing.T) {
	resp := SceneResponse{}
	typ := reflect.TypeOf(resp)

	fields := []string{"ID", "ChapterID", "SceneNumber", "Content", "CreatedAt", "UpdatedAt"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in SceneResponse", field)
		}
	}
}

func TestWorkflowRunResponse_HasFrozenCoreFields(t *testing.T) {
	resp := WorkflowRunResponse{}
	typ := reflect.TypeOf(resp)

	fields := []string{"ID", "ProjectID", "ChapterID", "Status", "CreatedAt", "UpdatedAt"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in WorkflowRunResponse", field)
		}
	}
}

func TestSystemCheckResultResponse_HasFrozenCoreFields(t *testing.T) {
	resp := SystemCheckResultResponse{}
	typ := reflect.TypeOf(resp)

	fields := []string{"Provider", "Severity", "Message"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in SystemCheckResultResponse", field)
		}
	}
}

func TestErrorResponse_HasFrozenFields(t *testing.T) {
	resp := ErrorResponse{}
	typ := reflect.TypeOf(resp)

	fields := []string{"Code", "Message", "Details"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in ErrorResponse", field)
		}
	}
}

func TestWebSocketEventEnvelope_HasFrozenFields(t *testing.T) {
	env := WebSocketEventEnvelope{}
	typ := reflect.TypeOf(env)

	fields := []string{"Type", "Timestamp", "Payload"}
	for _, field := range fields {
		if _, exists := typ.FieldByName(field); !exists {
			t.Fatalf("expected %s field in WebSocketEventEnvelope", field)
		}
	}
}

func TestWebSocketEventNames_AreFrozen(t *testing.T) {
	events := []string{
		EventWorkflowStarted,
		EventWorkflowStepChanged,
		EventWorkflowLog,
		EventWorkflowCompleted,
		EventWorkflowFailed,
		EventSystemCheckUpdated,
	}

	expected := []string{
		"workflow.started",
		"workflow.step_changed",
		"workflow.log",
		"workflow.completed",
		"workflow.failed",
		"system.check.updated",
	}

	for i, event := range events {
		if event != expected[i] {
			t.Fatalf("expected event name %q, got %q", expected[i], event)
		}
	}
}