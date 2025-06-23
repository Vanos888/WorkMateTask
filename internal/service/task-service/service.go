package task_service

import (
	io_bound_task "WorkMateTask/internal/jobs/task-processing/io-bound-task"
	"context"
	"github.com/google/uuid"
	"time"
)

type ITaskRepository interface {
	CreateTask(ctx context.Context, task io_bound_task.Task) error
	GetTaskStatus(ctx context.Context, taskID uuid.UUID) (io_bound_task.Task, error)
}

type IProcessingService interface {
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

type TaskService struct {
	tRepo             ITaskRepository
	processingService IProcessingService
}

func NewTaskService(tRepo ITaskRepository, processingService IProcessingService) *TaskService {
	return &TaskService{
		tRepo:             tRepo,
		processingService: processingService,
	}
}

func (t *TaskService) CreateTask(ctx context.Context, taskType io_bound_task.TaskType) (io_bound_task.Task, error) {

	task := io_bound_task.Task{
		Id:             uuid.New(),
		Type:           taskType,
		Status:         io_bound_task.TaskStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ProcessingTime: 0,
	}

	if err := t.tRepo.CreateTask(ctx, task); err != nil {
		return io_bound_task.Task{}, err
	}
	return task, nil

}
func (t *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return t.processingService.DeleteTask(ctx, id)
}

func (t *TaskService) GetTask(ctx context.Context, taskId uuid.UUID) (io_bound_task.Task, error) {
	task, err := t.tRepo.GetTaskStatus(ctx, taskId)
	if err != nil {
		return io_bound_task.Task{}, err
	}

	return task, nil

}
