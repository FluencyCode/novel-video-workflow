package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"novel-video-workflow/pkg/workflow"

	"github.com/gin-gonic/gin"
)

// WorkflowRunAPI handles workflow run endpoints.
type WorkflowRunAPI struct {
	executor *workflow.Executor
	storage  workflow.Storage
}

// NewWorkflowRunAPI creates a new workflow run API.
func NewWorkflowRunAPI(executor *workflow.Executor, storage workflow.Storage) *WorkflowRunAPI {
	return &WorkflowRunAPI{
		executor: executor,
		storage:  storage,
	}
}

// RegisterRoutes registers workflow run routes.
func (api *WorkflowRunAPI) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/workflow/runs", api.StartRun)
	router.GET("/api/workflow/runs/:id", api.GetRun)
}

// StartRun starts a new workflow run.
func (api *WorkflowRunAPI) StartRun(c *gin.Context) {
	var req struct {
		ChapterID uint `json:"chapter_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// TODO: Run asynchronously and broadcast progress via WebSocket
	result, err := api.executor.RunChapterWorkflow(context.Background(), workflow.RunRequest{
		ChapterID: req.ChapterID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "workflow_failed",
			Message: "Workflow execution failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Status: "success",
		Data:   result,
	})
}

// GetRun retrieves a workflow run by ID.
func (api *WorkflowRunAPI) GetRun(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_id",
			Message: "Invalid run ID",
		})
		return
	}

	run, err := api.storage.LoadByChapterID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "not_found",
			Message: "Workflow run not found",
		})
		return
	}

	c.JSON(http.StatusOK, WorkflowRunResponse{
		ID:           0, // TODO: Get from database
		ProjectID:    0, // TODO: Load from chapter
		ChapterID:    run.ChapterID,
		Status:       string(run.Status),
		CurrentStep:  string(run.CurrentStep),
		ErrorMessage: run.ErrorMessage,
		StartedAt:    derefTime(run.StartedAt),
		CompletedAt:  derefTime(run.FinishedAt),
		CreatedAt:    derefTime(run.StartedAt),
		UpdatedAt:    derefTime(run.FinishedAt),
	})
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}