package task

import (
	ogen_handler "WorkMateTask/internal/generated/servers/http/v1/task"
	io_bound_task "WorkMateTask/internal/jobs/task-processing/io-bound-task"
	"context"
	"github.com/google/uuid"
)

type ITaskService interface {
	CreateTask(ctx context.Context, taskType io_bound_task.TaskType) (io_bound_task.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
	GetTask(ctx context.Context, taskId uuid.UUID) (io_bound_task.Task, error)
}

type TaskApi struct {
	TaskService ITaskService
	ogen_handler.UnimplementedHandler
}

func NewTaskApi(taskService ITaskService) *TaskApi {
	return &TaskApi{
		TaskService: taskService,
	}
}
