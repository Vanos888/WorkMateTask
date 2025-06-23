package task

import (
	ogen_task "WorkMateTask/internal/generated/servers/http/v1/task"
	"context"
	"github.com/google/uuid"
	"time"
)

func (t *TaskApi) GetTaskStatus(ctx context.Context, req *ogen_task.GetTaskStatusRequest) (*ogen_task.GetTaskStatusResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	task, err := t.TaskService.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ogen_task.GetTaskStatusResponse{
		ID:        id.String(),
		Status:    task.Status.String(),
		CreatedAt: task.CreatedAt.UTC().Format(time.RFC3339),
		Duration:  time.Since(task.StartedAt).String(),
	}, nil
}
