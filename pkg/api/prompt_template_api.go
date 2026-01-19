package api

import (
	"net/http"
	"strconv"

	"novel-video-workflow/pkg/database"
	"novel-video-workflow/pkg/workflow"

	"github.com/gin-gonic/gin"
)

type PromptTemplateAPI struct {
	processor *workflow.Processor
}

func NewPromptTemplateAPI(processor *workflow.Processor) *PromptTemplateAPI {
	return &PromptTemplateAPI{
		processor: processor,
	}
}

// GetPromptTemplates godoc
// @Summary 获取所有激活的提示词模板
// @Description 获取所有激活的提示词模板列表
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Success 200 {array} database.PromptTemplate
// @Router /prompt-templates [get]
func (api *PromptTemplateAPI) GetPromptTemplates(c *gin.Context) {
	templates, err := database.GetActivePromptTemplates(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取提示词模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetPromptTemplateByID godoc
// @Summary 获取提示词模板详情
// @Description 根据ID获取提示词模板详情
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Success 200 {object} database.PromptTemplate
// @Router /prompt-templates/{id} [get]
func (api *PromptTemplateAPI) GetPromptTemplateByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	template, err := database.GetPromptTemplateByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetPromptTemplatesByCategory godoc
// @Summary 按分类获取提示词模板
// @Description 根据分类获取提示词模板列表
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Param category path string true "分类"
// @Success 200 {array} database.PromptTemplate
// @Router /prompt-templates/category/{category} [get]
func (api *PromptTemplateAPI) GetPromptTemplatesByCategory(c *gin.Context) {
	category := c.Param("category")

	templates, err := database.GetPromptTemplatesByCategory(database.DB, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取提示词模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// CreatePromptTemplate godoc
// @Summary 创建提示词模板
// @Description 创建一个新的提示词模板
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Param template body database.PromptTemplate true "模板信息"
// @Success 200 {object} database.PromptTemplate
// @Router /prompt-templates [post]
func (api *PromptTemplateAPI) CreatePromptTemplate(c *gin.Context) {
	var req database.PromptTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保用户创建的模板不是内置模板
	req.IsBuiltIn = false

	err := database.CreatePromptTemplate(database.DB, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建提示词模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// UpdatePromptTemplate godoc
// @Summary 更新提示词模板
// @Description 更新指定ID的提示词模板
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Param template body database.PromptTemplate true "模板更新信息"
// @Success 200 {object} database.PromptTemplate
// @Router /prompt-templates/{id} [put]
func (api *PromptTemplateAPI) UpdatePromptTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	var req database.PromptTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有模板
	existingTemplate, err := database.GetPromptTemplateByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 不允许修改内置模板
	if existingTemplate.IsBuiltIn {
		c.JSON(http.StatusForbidden, gin.H{"error": "不允许修改内置模板"})
		return
	}

	// 保持原有的一些属性不变
	req.ID = uint(id)
	req.IsBuiltIn = existingTemplate.IsBuiltIn // 保持原有内置状态
	req.CreatedAt = existingTemplate.CreatedAt // 保持创建时间

	err = database.UpdatePromptTemplate(database.DB, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新提示词模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeletePromptTemplate godoc
// @Summary 删除提示词模板
// @Description 删除指定ID的提示词模板
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Success 200 {object} map[string]string
// @Router /prompt-templates/{id} [delete]
func (api *PromptTemplateAPI) DeletePromptTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	template, err := database.GetPromptTemplateByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 不允许删除内置模板
	if template.IsBuiltIn {
		c.JSON(http.StatusForbidden, gin.H{"error": "不允许删除内置模板"})
		return
	}

	err = database.DeletePromptTemplate(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除提示词模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}