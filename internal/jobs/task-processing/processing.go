package task_processing

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type IJobBackground interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	GetID() uuid.UUID
}
type TaskProcessor struct {
	taskCh chan IJobBackground
	stopCh chan IJobBackground
}

func NewTaskProcessor(ctx context.Context, workerCount int) (*TaskProcessor, error) {
	taskCh := make(chan IJobBackground, workerCount)
	stopCh := make(chan IJobBackground, workerCount)
	proc := &TaskProcessor{
		taskCh: taskCh,
		stopCh: stopCh,
	}

	if err := proc.initBackgroundWorkers(ctx, workerCount); err != nil {
		return nil, err
	}

	return proc, nil
}

func (p *TaskProcessor) AddTask(task IJobBackground) {
	p.taskCh <- task
}

func (p *TaskProcessor) StopTask(task IJobBackground) {
	p.stopCh <- task
}

func (p *TaskProcessor) initBackgroundWorkers(ctx context.Context, workerCount int) error {
	// под каждый воркер создаем свой канал для отслеживания отмены исполнения таски
	cancelChs := make([]chan string, 0, workerCount)
	for i := 0; i < workerCount; i++ {
		cancelChs = append(cancelChs, make(chan string, 1))
	}

	// создаем воркеры и передаем каждому свой канал для отслеживания отмены таски
	for i := 0; i < workerCount; i++ {
		go p.worker(ctx, cancelChs[i])
	}

	// запускаем отслеживание отмены таски из вне
	go func() {
		for {
			select {
			case task := <-p.stopCh:
				// останавливаем таску
				if err := task.Stop(ctx); err != nil {
					fmt.Println(err)
				}

				// сигнализируем воркерам о том, что нужно отменить исполнение таски
				for _, ch := range cancelChs {
					ch <- task.GetID().String()
				}
			case <-ctx.Done():
				// если отменился главный контекст, то закрываем наши каналы
				close(p.stopCh)

				for _, ch := range cancelChs {
					close(ch)
				}

				close(p.taskCh)
				return
			}
		}
	}()

	return nil
}
func (p *TaskProcessor) worker(ctx context.Context, cancelCh chan string) {
	for task := range p.taskCh {
		// создаем отдельный контекст для таски, чтоб ее можно было остановить отдельно
		taskCtx, cancel := context.WithCancel(ctx)
		successCh := make(chan struct{}, 1)

		// наблюдаем за тем остановлена ли таска
		go func() {
			for {
				select {
				case taskID := <-cancelCh:
					if taskID == task.GetID().String() {
						cancel()
						return
					}
				case <-successCh:
					return
				}
			}
		}()

		if err := task.Run(taskCtx); err != nil {
			fmt.Println(err)
		}

		successCh <- struct{}{}
	}
}
