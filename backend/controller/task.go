package controller

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/client"
	"github.com/songquanpeng/go-api-starter/common/logger"
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
		"version": "1.0.0",
	})
}

type BatchTaskRequest struct {
	Files        []BatchFileInfo `json:"files" binding:"required"`
	OutputFormat string          `json:"output_format"`
	Options      ParseOptions    `json:"options"`
}

type BatchFileInfo struct {
	FileName string `json:"file_name" binding:"required"`
	FileSize int64  `json:"file_size"`
	FileURL  string `json:"file_url"`
}

type ParseOptions struct {
	PreserveTables   bool   `json:"preserve_tables"`
	ExtractImages    bool   `json:"extract_images"`
	LatexFormulas    bool   `json:"latex_formulas"`
	PreserveHeadings bool   `json:"preserve_headings"`
	Language         string `json:"language"`
}

type BatchTaskResponse struct {
	TaskIDs   []string `json:"task_ids"`
	Total     int      `json:"total"`
	CreatedAt string   `json:"created_at"`
}

func (h *TaskController) CreateTask(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.ParamError(c, "请上传文件")
		return
	}

	outputFormat := c.DefaultPostForm("output_format", "markdown")
	if !isValidFormat(outputFormat) {
		common.ParamError(c, "输出格式必须是 markdown、txt 或 json")
		return
	}

	taskID := generateTaskID()
	ext := filepath.Ext(file.Filename)
	savePath := filepath.Join(common.GlobalConfig.Upload.StoragePath, taskID+ext)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		logger.Error("保存文件失败:", err)
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
		Progress:     0,
	}

	if userID, exists := getUserID(c); exists {
		task.UserID = userID
	}

	if err := common.DB.Create(task).Error; err != nil {
		logger.Error("创建任务失败:", err)
		common.ServerError(c, "创建任务失败")
		return
	}

	go h.processTaskAsync(taskID)

	common.Success(c, gin.H{
		"task_id":    taskID,
		"file_name":  file.Filename,
		"file_size":  file.Size,
		"status":     task.Status,
		"progress":   0,
		"created_at": task.CreatedAt.Format(time.RFC3339),
	})
}

func (h *TaskController) CreateBatchTasks(c *gin.Context) {
	var req BatchTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c, "请求参数格式错误")
		return
	}

	if len(req.Files) == 0 {
		common.ParamError(c, "请至少上传一个文件")
		return
	}

	if len(req.Files) > 50 {
		common.ParamError(c, "单次最多上传50个文件")
		return
	}

	outputFormat := req.OutputFormat
	if outputFormat == "" {
		outputFormat = "markdown"
	}
	if !isValidFormat(outputFormat) {
		common.ParamError(c, "输出格式必须是 markdown、txt 或 json")
		return
	}

	userID, hasUser := getUserID(c)
	taskIDs := make([]string, 0, len(req.Files))

	for _, fileInfo := range req.Files {
		taskID := generateTaskID()

		task := &model.ParseTask{
			TaskID:       taskID,
			FileName:     fileInfo.FileName,
			FilePath:     "",
			FileSize:     fileInfo.FileSize,
			OutputFormat: outputFormat,
			Status:       model.TaskStatusPending,
			Progress:     0,
		}

		if hasUser {
			task.UserID = userID
		}

		if err := common.DB.Create(task).Error; err != nil {
			logger.Error("创建批量任务失败:", err)
			continue
		}

		taskIDs = append(taskIDs, taskID)
	}

	for _, taskID := range taskIDs {
		go h.processTaskAsync(taskID)
	}

	common.Success(c, BatchTaskResponse{
		TaskIDs:   taskIDs,
		Total:     len(taskIDs),
		CreatedAt: time.Now().Format(time.RFC3339),
	})
}

func (h *TaskController) UploadBatchFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		common.ParamError(c, "请上传文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		common.ParamError(c, "请至少上传一个文件")
		return
	}

	if len(files) > 50 {
		common.ParamError(c, "单次最多上传50个文件")
		return
	}

	outputFormat := c.DefaultPostForm("output_format", "markdown")
	if !isValidFormat(outputFormat) {
		outputFormat = "markdown"
	}

	userID, hasUser := getUserID(c)
	taskIDs := make([]string, 0, len(files))
	tasks := make([]gin.H, 0, len(files))

	for _, file := range files {
		taskID := generateTaskID()
		ext := filepath.Ext(file.Filename)
		savePath := filepath.Join(common.GlobalConfig.Upload.StoragePath, taskID+ext)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			logger.Error("保存文件失败:", err, "filename:", file.Filename)
			continue
		}

		task := &model.ParseTask{
			TaskID:       taskID,
			FileName:     file.Filename,
			FilePath:     savePath,
			FileSize:     file.Size,
			OutputFormat: outputFormat,
			Status:       model.TaskStatusPending,
			Progress:     0,
		}

		if hasUser {
			task.UserID = userID
		}

		if err := common.DB.Create(task).Error; err != nil {
			logger.Error("创建任务失败:", err)
			continue
		}

		taskIDs = append(taskIDs, taskID)
		tasks = append(tasks, gin.H{
			"task_id":    taskID,
			"file_name":  file.Filename,
			"file_size":  file.Size,
			"status":     task.Status,
			"progress":   0,
			"created_at": task.CreatedAt.Format(time.RFC3339),
		})
	}

	for _, taskID := range taskIDs {
		go h.processTaskAsync(taskID)
	}

	common.Success(c, gin.H{
		"task_ids":   taskIDs,
		"tasks":      tasks,
		"total":      len(taskIDs),
		"created_at": time.Now().Format(time.RFC3339),
	})
}

func (h *TaskController) GetBatchStatus(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		common.ParamError(c, "请提供任务ID列表")
		return
	}

	taskIDs := strings.Split(idsStr, ",")
	if len(taskIDs) > 100 {
		common.ParamError(c, "单次最多查询100个任务")
		return
	}

	var tasks []model.ParseTask
	common.DB.Where("task_id IN ?", taskIDs).Find(&tasks)

	taskMap := make(map[string]gin.H)
	for _, task := range tasks {
		taskMap[task.TaskID] = gin.H{
			"task_id":    task.TaskID,
			"file_name":  task.FileName,
			"status":     task.Status,
			"progress":   task.Progress,
			"error_msg":  task.ErrorMsg,
			"updated_at": task.UpdatedAt.Format(time.RFC3339),
		}
	}

	results := make([]gin.H, 0, len(taskIDs))
	for _, id := range taskIDs {
		if task, ok := taskMap[id]; ok {
			results = append(results, task)
		} else {
			results = append(results, gin.H{
				"task_id": id,
				"status":  "not_found",
			})
		}
	}

	common.Success(c, gin.H{
		"tasks": results,
		"total": len(results),
	})
}

func (h *TaskController) DownloadBatchResults(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		common.ParamError(c, "请提供任务ID列表")
		return
	}

	taskIDs := strings.Split(idsStr, ",")
	if len(taskIDs) > 50 {
		common.ParamError(c, "单次最多下载50个任务结果")
		return
	}

	var results []model.ParseResult
	common.DB.Where("task_id IN ?", taskIDs).Find(&results)

	if len(results) == 0 {
		common.NotFound(c, "没有找到可下载的结果")
		return
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, result := range results {
		var task model.ParseTask
		if err := common.DB.Where("task_id = ?", result.TaskID).First(&task).Error; err != nil {
			continue
		}

		ext := ".md"
		if task.OutputFormat == "txt" {
			ext = ".txt"
		} else if task.OutputFormat == "json" {
			ext = ".json"
		}

		fileName := task.FileName
		if idx := strings.LastIndex(fileName, "."); idx > 0 {
			fileName = fileName[:idx]
		}
		fileName = fileName + ext

		w, err := zipWriter.Create(fileName)
		if err != nil {
			continue
		}
		w.Write([]byte(result.Content))
	}

	zipWriter.Close()

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=results_%s.zip", time.Now().Format("20060102150405")))
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

func (h *TaskController) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var tasks []model.ParseTask
	var total int64

	query := common.DB.Model(&model.ParseTask{})

	if userID, exists := getUserID(c); exists {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("file_name LIKE ?", "%"+search+"%")
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks)

	stats := h.getTaskStats(c)

	common.Success(c, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"stats": stats,
	})
}

func (h *TaskController) GetTaskStats(c *gin.Context) {
	stats := h.getTaskStats(c)
	common.Success(c, stats)
}

func (h *TaskController) getTaskStats(c *gin.Context) gin.H {
	query := common.DB.Model(&model.ParseTask{})

	if userID, exists := getUserID(c); exists {
		query = query.Where("user_id = ?", userID)
	}

	var total, pending, processing, completed, failed int64
	query.Count(&total)
	query.Where("status = ?", model.TaskStatusPending).Count(&pending)
	query.Where("status = ?", model.TaskStatusProcessing).Count(&processing)
	query.Where("status = ?", model.TaskStatusCompleted).Count(&completed)
	query.Where("status = ?", model.TaskStatusFailed).Count(&failed)

	return gin.H{
		"total":      total,
		"pending":    pending,
		"processing": processing,
		"completed":  completed,
		"failed":     failed,
	}
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
		"task_id":     task.TaskID,
		"file_name":   task.FileName,
		"status":      task.Status,
		"progress":    task.Progress,
		"error_msg":   task.ErrorMsg,
		"output_path": task.FilePath,
		"created_at":  task.CreatedAt.Format(time.RFC3339),
		"updated_at":  task.UpdatedAt.Format(time.RFC3339),
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

	common.Success(c, gin.H{
		"task_id":    taskID,
		"file_name":  task.FileName,
		"content":    result.Content,
		"format":     task.OutputFormat,
		"created_at": result.CreatedAt.Format(time.RFC3339),
	})
}

func (h *TaskController) DownloadResult(c *gin.Context) {
	taskID := c.Param("task_id")
	format := c.DefaultQuery("format", "md")

	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		common.NotFound(c, "任务不存在")
		return
	}

	var result model.ParseResult
	if err := common.DB.Where("task_id = ?", taskID).First(&result).Error; err != nil {
		common.NotFound(c, "解析结果不存在")
		return
	}

	fileName := task.FileName
	if idx := strings.LastIndex(fileName, "."); idx > 0 {
		fileName = fileName[:idx]
	}

	ext := "." + format
	if format == "md" {
		ext = ".md"
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s%s", fileName, ext))
	c.String(http.StatusOK, result.Content)
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

func (h *TaskController) DeleteBatchTasks(c *gin.Context) {
	var req struct {
		TaskIDs []string `json:"task_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c, "请求参数格式错误")
		return
	}

	if len(req.TaskIDs) == 0 {
		common.ParamError(c, "请提供要删除的任务ID")
		return
	}

	if len(req.TaskIDs) > 100 {
		common.ParamError(c, "单次最多删除100个任务")
		return
	}

	result := common.DB.Where("task_id IN ?", req.TaskIDs).Delete(&model.ParseTask{})
	common.DB.Where("task_id IN ?", req.TaskIDs).Delete(&model.ParseResult{})

	common.Success(c, gin.H{
		"deleted": result.RowsAffected,
	})
}

func (h *TaskController) processTaskAsync(taskID string) {
	var task model.ParseTask
	if err := common.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		logger.Error("任务不存在:", taskID)
		return
	}

	h.updateTaskProgress(taskID, model.TaskStatusProcessing, 10)

	minoruClient := client.NewMineruClient(&common.GlobalConfig.Minoru)

	h.updateTaskProgress(taskID, model.TaskStatusProcessing, 30)

	content, err := minoruClient.Parse(task.FilePath, task.OutputFormat)
	if err != nil {
		h.updateTaskError(taskID, err.Error())
		return
	}

	h.updateTaskProgress(taskID, model.TaskStatusProcessing, 80)

	resultFilePath := task.FilePath + "." + task.OutputFormat
	if err := common.WriteFile(resultFilePath, []byte(content)); err != nil {
		h.updateTaskError(taskID, err.Error())
		return
	}

	result := &model.ParseResult{
		TaskID:   taskID,
		Content:  content,
		FilePath: resultFilePath,
	}
	common.DB.Create(result)

	h.updateTaskProgress(taskID, model.TaskStatusCompleted, 100)

	h.notifyTaskComplete(taskID, task.FileName)
}

func (h *TaskController) updateTaskProgress(taskID string, status string, progress int) {
	common.DB.Model(&model.ParseTask{}).Where("task_id = ?", taskID).Updates(map[string]interface{}{
		"status":   status,
		"progress": progress,
	})
}

func (h *TaskController) updateTaskError(taskID string, errMsg string) {
	common.DB.Model(&model.ParseTask{}).Where("task_id = ?", taskID).Updates(map[string]interface{}{
		"status":     model.TaskStatusFailed,
		"error_msg":  errMsg,
		"progress":   0,
	})
}

func (h *TaskController) notifyTaskComplete(taskID string, fileName string) {
	if common.RDB != nil {
		notifyData, _ := json.Marshal(map[string]interface{}{
			"type":      "task_complete",
			"task_id":   taskID,
			"file_name": fileName,
			"timestamp": time.Now().Unix(),
		})
		common.RDB.Publish(nil, "task:notifications", string(notifyData))
	}
}

func generateTaskID() string {
	return "task_" + uuid.New().String()[:24]
}

func isValidFormat(format string) bool {
	return format == "markdown" || format == "txt" || format == "json"
}

func getUserID(c *gin.Context) (string, bool) {
	user, exists := c.Get("user")
	if !exists {
		return "", false
	}
	u, ok := user.(*model.User)
	if !ok {
		return "", false
	}
	return u.UserID, true
}

func (h *TaskController) UploadFileURL(c *gin.Context) {
	var req struct {
		URL          string `json:"url" binding:"required"`
		FileName     string `json:"file_name"`
		OutputFormat string `json:"output_format"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c, "请求参数格式错误")
		return
	}

	outputFormat := req.OutputFormat
	if outputFormat == "" {
		outputFormat = "markdown"
	}
	if !isValidFormat(outputFormat) {
		common.ParamError(c, "输出格式必须是 markdown、txt 或 json")
		return
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		common.ParamError(c, "无法访问文件URL")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.ParamError(c, "文件下载失败")
		return
	}

	taskID := generateTaskID()
	fileName := req.FileName
	if fileName == "" {
		fileName = filepath.Base(req.URL)
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".pdf"
		fileName = fileName + ext
	}

	savePath := filepath.Join(common.GlobalConfig.Upload.StoragePath, taskID+ext)

	file, err := os.Create(savePath)
	if err != nil {
		common.ServerError(c, "保存文件失败")
		return
	}
	defer file.Close()

	fileSize, err := io.Copy(file, resp.Body)
	if err != nil {
		common.ServerError(c, "保存文件失败")
		return
	}

	task := &model.ParseTask{
		TaskID:       taskID,
		FileName:     fileName,
		FilePath:     savePath,
		FileSize:     fileSize,
		OutputFormat: outputFormat,
		Status:       model.TaskStatusPending,
		Progress:     0,
	}

	if userID, exists := getUserID(c); exists {
		task.UserID = userID
	}

	if err := common.DB.Create(task).Error; err != nil {
		common.ServerError(c, "创建任务失败")
		return
	}

	go h.processTaskAsync(taskID)

	common.Success(c, gin.H{
		"task_id":    taskID,
		"file_name":  fileName,
		"file_size":  fileSize,
		"status":     task.Status,
		"progress":   0,
		"created_at": task.CreatedAt.Format(time.RFC3339),
	})
}
