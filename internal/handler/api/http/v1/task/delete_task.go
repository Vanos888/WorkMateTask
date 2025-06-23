package task

import (
	"WorkMateTask/internal/generated/servers/http/v1/task"
	"context"
	"github.com/google/uuid"
)

func (t *TaskApi) DeleteTask(ctx context.Context, req *task.DeleteTaskRequest) (*task.DeleteTaskResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	if err = t.TaskService.DeleteTask(ctx, id); err != nil {
		return nil, err
	}
	return &task.DeleteTaskResponse{}, nil
}
