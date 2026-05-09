package handler

import (
	"pdf-parser/internal/model"
	"pdf-parser/internal/service"
	"pdf-parser/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ParamError(c, "请上传PDF文件")
		return
	}

	outputFormat := c.DefaultPostForm("output_format", "markdown")
	if outputFormat != "markdown" && outputFormat != "txt" {
		response.ParamError(c, "输出格式必须是 markdown 或 txt")
		return
	}

	result, err := h.taskService.CreateTask(c.Request.Context(), file, model.OutputFormat(outputFormat))
	if err != nil {
		if err == model.ErrInvalidFileType {
			response.ParamError(c, err.Error())
		} else if err == model.ErrFileTooLarge {
			response.ParamError(c, err.Error())
		} else {
			response.ServerError(c, "创建任务失败")
		}
		return
	}

	response.Success(c, result)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

	status := c.Query("status")

	result, err := h.taskService.ListTasks(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.ServerError(c, "获取任务列表失败")
		return
	}

	response.Success(c, result)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.ParamError(c, "任务ID不能为空")
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), taskID)
	if err != nil {
		if err == model.ErrTaskNotFound {
			response.NotFound(c, "任务不存在")
		} else {
			response.ServerError(c, "获取任务详情失败")
		}
		return
	}

	response.Success(c, task)
}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.ParamError(c, "任务ID不能为空")
		return
	}

	status, err := h.taskService.GetTaskStatus(c.Request.Context(), taskID)
	if err != nil {
		if err == model.ErrTaskNotFound {
			response.NotFound(c, "任务不存在")
		} else {
			response.ServerError(c, "获取任务状态失败")
		}
		return
	}

	response.Success(c, status)
}

func (h *TaskHandler) GetTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.ParamError(c, "任务ID不能为空")
		return
	}

	result, err := h.taskService.GetTaskResult(c.Request.Context(), taskID)
	if err != nil {
		if err == model.ErrTaskNotFound {
			response.NotFound(c, "任务不存在")
		} else if err == model.ErrResultNotFound {
			response.NotFound(c, "解析结果不存在")
		} else {
			response.ServerError(c, "获取解析结果失败")
		}
		return
	}

	response.Success(c, result)
}

func (h *TaskHandler) DownloadResult(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.ParamError(c, "任务ID不能为空")
		return
	}

	err := h.taskService.DownloadResult(c.Request.Context(), taskID, c)
	if err != nil {
		if err == model.ErrTaskNotFound {
			response.NotFound(c, "任务不存在")
		} else if err == model.ErrResultNotFound {
			response.NotFound(c, "解析结果不存在")
		} else {
			response.ServerError(c, "下载失败")
		}
		return
	}
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.ParamError(c, "任务ID不能为空")
		return
	}

	err := h.taskService.DeleteTask(c.Request.Context(), taskID)
	if err != nil {
		if err == model.ErrTaskNotFound {
			response.NotFound(c, "任务不存在")
		} else {
			response.ServerError(c, "删除任务失败")
		}
		return
	}

	response.SuccessWithMessage(c, "任务删除成功")
}

func (h *TaskHandler) HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
		"service": "pdf-parser",
	})
}
