package task_repository

import (
	io_bound_task "WorkMateTask/internal/jobs/task-processing/io-bound-task"
	"WorkMateTask/internal/storage"
	"context"
	"errors"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("task not found")

type TaskRepository struct {
	tStorage *storage.Storage
}

func NewTaskRepository(tStorage *storage.Storage) *TaskRepository {
	return &TaskRepository{
		tStorage: tStorage,
	}
}

func (t *TaskRepository) CreateTask(ctx context.Context, task io_bound_task.Task) error {
	t.tStorage.SetTask(task.Id.String(), task)
	return nil
}

func (t *TaskRepository) GetTaskStatus(ctx context.Context, taskID uuid.UUID) (io_bound_task.Task, error) {
	task, ok := t.tStorage.GetTask(taskID.String())
	if !ok {
		return io_bound_task.Task{}, ErrNotFound
	}
	return task, nil

}
