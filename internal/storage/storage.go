package storage

import (
	io_bound_task "WorkMateTask/internal/jobs/task-processing/io-bound-task"
	"sync"
)

type Storage struct {
	taskStorage map[string]io_bound_task.Task
	mu          sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		taskStorage: make(map[string]io_bound_task.Task),
	}
}

func (s *Storage) GetTask(key string) (io_bound_task.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.taskStorage[key]
	return val, ok
}

func (s *Storage) SetTask(key string, value io_bound_task.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.taskStorage[key] = value
}
func (s *Storage) DeleteTask(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.taskStorage, key)
}
func (s *Storage) GetTasksByStatus(status io_bound_task.TaskStatus) ([]io_bound_task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]io_bound_task.Task, 0, 10)

	for _, task := range s.taskStorage {
		if task.Status == status {
			tasks = append(tasks, task)
		}
		if len(tasks) == 10 {
			return tasks, nil
		}
	}
	return tasks, nil
}
