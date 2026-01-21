package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"novel-video-workflow/pkg/api"
	"novel-video-workflow/pkg/database"
	"novel-video-workflow/pkg/workflow"

	"go.uber.org/zap"
)

func main() {
	// 创建logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 初始化数据库
	dbPath := filepath.Join(os.TempDir(), "test_chapter_workflow.db")
	err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建工作流处理器
	processor, err := workflow.NewProcessor(logger)
	if err != nil {
		log.Fatalf("创建工作流处理器失败: %v", err)
	}

	// 创建章节工作流API
	chapterWorkflowAPI := api.NewChapterWorkflowAPI(processor)

	// 测试创建项目
	project, err := processor.CreateProject("测试项目", "这是一个测试项目", "测试提示词", "test123")
	if err != nil {
		log.Fatalf("创建项目失败: %v", err)
	}
	fmt.Printf("创建项目成功: ID=%d, 名称=%s\n", project.ID, project.Name)

	// 测试创建章节
	chapter, err := processor.CreateChapter(project.ID, "第一章", "这是第一章的内容", "第一章氛围提示词")
	if err != nil {
		log.Fatalf("创建章节失败: %v", err)
	}
	fmt.Printf("创建章节成功: ID=%d, 标题=%s\n", chapter.ID, chapter.Title)

	// 测试创建工作流步骤
	step, err := processor.CreateWorkflowStep(chapter.ID, "音频生成", "pending", `{"input": "test"}`, 1)
	if err != nil {
		log.Fatalf("创建工作流步骤失败: %v", err)
	}
	fmt.Printf("创建工作流步骤成功: ID=%d, 名称=%s, 状态=%s\n", step.ID, step.Name, step.Status)

	// 测试获取章节的工作流步骤
	steps, err := processor.GetWorkflowStepsByChapterID(chapter.ID)
	if err != nil {
		log.Fatalf("获取工作流步骤失败: %v", err)
	}
	fmt.Printf("获取工作流步骤成功: 共%d个步骤\n", len(steps))

	// 测试更新工作流步骤
	err = processor.UpdateWorkflowStep(step.ID, map[string]interface{}{
		"status": "completed",
		"error":  "",
	})
	if err != nil {
		log.Fatalf("更新工作流步骤失败: %v", err)
	}
	fmt.Printf("更新工作流步骤成功: ID=%d, 状态变为completed\n", step.ID)

	// 测试更新章节状态
	err = processor.UpdateChapterStatusBasedOnSteps(chapter.ID)
	if err != nil {
		log.Fatalf("更新章节状态失败: %v", err)
	}
	fmt.Printf("更新章节状态成功\n")

	// 获取更新后的章节
	updatedChapter, err := processor.GetChapterByID(chapter.ID)
	if err != nil {
		log.Fatalf("获取章节失败: %v", err)
	}
	fmt.Printf("更新后的章节状态: %s\n", updatedChapter.Status)

	// 测试重试工作流步骤
	err = processor.UpdateWorkflowStep(step.ID, map[string]interface{}{
		"status":      "pending",
		"retry_count": step.RetryCount + 1,
	})
	if err != nil {
		log.Fatalf("重试工作流步骤失败: %v", err)
	}
	fmt.Printf("重试工作流步骤成功: ID=%d, 重试次数增加\n", step.ID)

	// 测试获取章节的所有信息
	chapters, err := processor.GetChaptersByProjectID(project.ID)
	if err != nil {
		log.Fatalf("获取项目章节失败: %v", err)
	}
	fmt.Printf("获取项目章节成功: 共%d个章节\n", len(chapters))

	// 测试删除工作流步骤
	err = processor.DeleteWorkflowStep(step.ID)
	if err != nil {
		log.Fatalf("删除工作流步骤失败: %v", err)
	}
	fmt.Printf("删除工作流步骤成功: ID=%d\n", step.ID)

	// 测试删除章节
	err = processor.DeleteChapter(chapter.ID)
	if err != nil {
		log.Fatalf("删除章节失败: %v", err)
	}
	fmt.Printf("删除章节成功: ID=%d\n", chapter.ID)

	// 测试删除项目
	err = processor.DeleteProject(project.ID)
	if err != nil {
		log.Fatalf("删除项目失败: %v", err)
	}
	fmt.Printf("删除项目成功: ID=%d\n", project.ID)

	fmt.Println("所有测试通过！章节工作流功能正常工作。")
}