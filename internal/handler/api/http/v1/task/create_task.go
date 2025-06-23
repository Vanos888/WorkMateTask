package task

import (
	"WorkMateTask/internal/generated/servers/http/v1/task"
	io_bound_task "WorkMateTask/internal/jobs/task-processing/io-bound-task"
	"context"
)

func (t *TaskApi) CreateTask(ctx context.Context, req *task.CreateTaskRequest) (*task.CreateTaskResponse, error) {
	taskType := io_bound_task.StringToTaskType(req.Type)
	createdTask, err := t.TaskService.CreateTask(ctx, taskType)
	if err != nil {
		return nil, err
	}

	return &task.CreateTaskResponse{
		ID: createdTask.Id.String(),
	}, nil

}
