package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"novel-video-workflow/pkg/workflow"

	"github.com/gin-gonic/gin"
)

// WorkflowTrackingAPI 处理工作流跟踪相关的API
type WorkflowTrackingAPI struct {
	processor *workflow.Processor
}

// NewWorkflowTrackingAPI 创建新的工作流跟踪API实例
func NewWorkflowTrackingAPI(processor *workflow.Processor) *WorkflowTrackingAPI {
	return &WorkflowTrackingAPI{
		processor: processor,
	}
}

// RecordChapterWorkflowParams 记录章节工作流参数
func (api *WorkflowTrackingAPI) RecordChapterWorkflowParams(c *gin.Context) {
	var req struct {
		ChapterID      uint                   `json:"chapter_id" binding:"required"`
		WorkflowParams map[string]interface{} `json:"workflow_params" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 序列化工作流参数
	paramsBytes, err := json.Marshal(req.WorkflowParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化工作流参数失败: " + err.Error()})
		return
	}

	paramsStr := string(paramsBytes)

	// 更新章节的工作流参数
	err = api.processor.UpdateChapter(req.ChapterID, map[string]interface{}{
		"workflow_params": paramsStr,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新章节工作流参数失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "章节工作流参数记录成功",
		"chapter_id": req.ChapterID,
	})
}

// GetChapterWorkflowParams 获取章节工作流参数
func (api *WorkflowTrackingAPI) GetChapterWorkflowParams(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节ID"})
		return
	}

	chapter, err := api.processor.GetChapterByID(uint(chapterID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
		return
	}

	var workflowParams map[string]interface{}
	if chapter.WorkflowParams != "" {
		err = json.Unmarshal([]byte(chapter.WorkflowParams), &workflowParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析工作流参数失败: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"chapter_id":       chapter.ID,
		"workflow_params":  workflowParams,
		"title":            chapter.Title,
		"updated_at":       chapter.UpdatedAt,
	})
}

// RecordSceneWorkflowParams 记录场景工作流参数
func (api *WorkflowTrackingAPI) RecordSceneWorkflowParams(c *gin.Context) {
	var req struct {
		SceneID          uint                   `json:"scene_id" binding:"required"`
		WorkflowDetails  map[string]interface{} `json:"workflow_details" binding:"required"`
		Status           string                 `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"workflow_details": req.WorkflowDetails,
	}

	if req.Status != "" {
		updates["status"] = req.Status
	}

	// 更新场景的工作流参数
	err := api.processor.UpdateSceneWithWorkflowDetails(req.SceneID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新场景工作流参数失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "场景工作流参数记录成功",
		"scene_id": req.SceneID,
		"status":   req.Status,
	})
}

// GetSceneWorkflowParams 获取场景工作流参数
func (api *WorkflowTrackingAPI) GetSceneWorkflowParams(c *gin.Context) {
	sceneID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的场景ID"})
		return
	}

	scene, err := api.processor.GetSceneByID(uint(sceneID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "场景不存在"})
		return
	}

	var workflowDetails map[string]interface{}
	if scene.WorkflowDetails != "" {
		err = json.Unmarshal([]byte(scene.WorkflowDetails), &workflowDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析工作流详细参数失败: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"scene_id":         scene.ID,
		"workflow_details": workflowDetails,
		"status":           scene.Status,
		"title":            scene.Title,
		"start_time":       scene.StartTime,
		"end_time":         scene.EndTime,
		"updated_at":       scene.UpdatedAt,
	})
}

// GetScenesByChapter 获取章节下的所有场景及工作流参数
func (api *WorkflowTrackingAPI) GetScenesByChapter(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节ID"})
		return
	}

	scenes, err := api.processor.GetScenesByChapterID(uint(chapterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取场景列表失败: " + err.Error()})
		return
	}

	var processedScenes []map[string]interface{}
	for _, scene := range scenes {
		var workflowDetails map[string]interface{}
		if scene.WorkflowDetails != "" {
			err = json.Unmarshal([]byte(scene.WorkflowDetails), &workflowDetails)
			if err != nil {
				workflowDetails = map[string]interface{}{"error": "解析工作流详细参数失败"}
			}
		}

		processedScene := map[string]interface{}{
			"id":               scene.ID,
			"title":            scene.Title,
			"description":      scene.Description,
			"prompt":           scene.Prompt,
			"status":           scene.Status,
			"start_time":       scene.StartTime,
			"end_time":         scene.EndTime,
			"workflow_details": workflowDetails,
			"created_at":       scene.CreatedAt,
			"updated_at":       scene.UpdatedAt,
		}

		processedScenes = append(processedScenes, processedScene)
	}

	c.JSON(http.StatusOK, gin.H{
		"chapter_id": uint(chapterID),
		"scenes":     processedScenes,
	})
}