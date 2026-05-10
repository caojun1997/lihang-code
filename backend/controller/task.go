package controller

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/logger"
	"github.com/songquanpeng/go-api-starter/common/client"
	"github.com/songquanpeng/go-api-starter/model"
)

type TaskController struct{}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (h *TaskController) HealthCheck(c *gin.Context) {
	common.Success(c, gin.H{
		"status":  "ok",
		"service": "pdf-parser",
	})
}

func (h *TaskController) CreateTask(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.ParamError(c, "请上传PDF文件")
		return
	}

	outputFormat := c.DefaultPostForm("output_format", "markdown")
	if outputFormat != "markdown" && outputFormat != "txt" {
		common.ParamError(c, "输出格式必须是 markdown 或 txt")
		return
	}

	taskID := common.GetRandomCode(32)
	ext := filepath.Ext(file.Filename)
	savePath := filepath.Join(common.GlobalConfig.Upload.StoragePath, taskID+ext)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		logger.Error("Failed to save uploaded file:", err)
		common.ServerError(c, "保存文件失败")
		return
	}

	task := &model.ParseTask{
		TaskID:       taskID,
		FileName:     file.Filename,
		FilePath:     savePath,
		FileSize:     file.Size,
		OutputFormat: outputFormat,
		Status:       model.TaskStatusPending,
	}

	user, exists := c.Get("user")
	if exists {
		task.UserID = user.(*model.User).UserID
	}

	if err := common.DB.Create(task).Error; err != nil {
		logger.Error("Failed to create task:", err)
		common.ServerError(c, "创建任务失败")
		return
	}

	go h.processTask(taskID)

	common.Success(c, gin.H{
		"task_id":    taskID,
		"file_name":  file.Filename,
		"file_size":  file.Size,
		"status":     task.Status,
		"created_at": task.CreatedAt,
	})
}

func (h *TaskController) processTask(taskID string) {
	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		logger.Error("Failed to find task:", err)
		return
	}

	common.DB.Model(&task).Update("status", model.TaskStatusProcessing)

	minoruClient := client.NewMineruClient(&common.GlobalConfig.Minoru)
	content, err := minoruClient.Parse(task.FilePath, task.OutputFormat)
	if err != nil {
		common.DB.Model(&task).Updates(map[string]interface{}{
			"status":    model.TaskStatusFailed,
			"error_msg": err.Error(),
		})
		return
	}

	resultFilePath := task.FilePath + "." + task.OutputFormat
	if err := common.WriteFile(resultFilePath, []byte(content)); err != nil {
		common.DB.Model(&task).Updates(map[string]interface{}{
			"status":    model.TaskStatusFailed,
			"error_msg": err.Error(),
		})
		return
	}

	result := &model.ParseResult{
		TaskID:   taskID,
		Content:  content,
		FilePath: resultFilePath,
	}
	common.DB.Create(result)

	common.DB.Model(&task).Update("status", model.TaskStatusCompleted)
}

func (h *TaskController) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var tasks []model.ParseTask
	var total int64

	query := common.DB.Model(&model.ParseTask{})

	user, exists := c.Get("user")
	if exists {
		query = query.Where("user_id = ?", user.(*model.User).UserID)
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks)

	common.Success(c, gin.H{
		"tasks": tasks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *TaskController) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")

	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		common.NotFound(c, "任务不存在")
		return
	}

	common.Success(c, task)
}

func (h *TaskController) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")

	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		common.NotFound(c, "任务不存在")
		return
	}

	common.Success(c, gin.H{
		"task_id": task.TaskID,
		"status":  task.Status,
	})
}

func (h *TaskController) GetTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")

	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		common.NotFound(c, "任务不存在")
		return
	}

	if task.Status != model.TaskStatusCompleted {
		common.ParamError(c, "任务还未完成")
		return
	}

	var result model.ParseResult
	if err := common.DB.Where("task_id = ?", taskID).First(&result).Error; err != nil {
		common.NotFound(c, "解析结果不存在")
		return
	}

	common.Success(c, result)
}

func (h *TaskController) DownloadResult(c *gin.Context) {
	taskID := c.Param("task_id")

	var result model.ParseResult
	if err := common.DB.Where("task_id = ?", taskID).First(&result).Error; err != nil {
		common.NotFound(c, "解析结果不存在")
		return
	}

	c.File(result.FilePath)
}

func (h *TaskController) DeleteTask(c *gin.Context) {
	taskID := c.Param("task_id")

	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		common.NotFound(c, "任务不存在")
		return
	}

	common.DB.Delete(&task)
	common.DB.Where("task_id = ?", taskID).Delete(&model.ParseResult{})

	common.SuccessWithMessage(c, "任务删除成功")
}
