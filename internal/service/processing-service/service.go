package processing_service

import (
	task_processing "WorkMateTask/internal/jobs/task-processing"
	io_bound_task "WorkMateTask/internal/jobs/task-processing/io-bound-task"
	"WorkMateTask/internal/storage"
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"time"
)

type ProcessingService struct {
	processing IProcessing
	storage    *storage.Storage
}

type IProcessing interface {
	AddTask(task task_processing.IJobBackground)
	StopTask(task task_processing.IJobBackground)
}

func NewTaskProcessingService(processing IProcessing, storage *storage.Storage) *ProcessingService {
	return &ProcessingService{
		processing: processing,
		storage:    storage,
	}
}

func (p *ProcessingService) ObserveTasks(ctx context.Context) error {
	//получает задачи со статусом пендинг из сторейджа и кладет их в процессинг через AddTask + вечный цикл
	//меняет статус задачи в торейдже через репозиторий
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tasks, err := p.storage.GetTasksByStatus(io_bound_task.TaskStatusPending)
				if err != nil {
					fmt.Println("Error while observing tasks: ", err)
				}
				for _, task := range tasks {
					p.processing.AddTask(&task)
					task.Status = io_bound_task.TaskStatusActive
					task.StartedAt = time.Now()
					p.storage.SetTask(task.Id.String(), task)
				}
			}

		}
	}()
	return nil
}

func (p *ProcessingService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	//вызывает метод стоп в процессинге и удаляет задачу из сторейджа через репозиторий
	task, ok := p.storage.GetTask(id.String())
	if !ok {
		return errors.New("Task not found")
	}
	if task.Status == io_bound_task.TaskStatusActive {
		p.processing.StopTask(&task)
	}
	p.storage.DeleteTask(task.Id.String())
	return nil
}
