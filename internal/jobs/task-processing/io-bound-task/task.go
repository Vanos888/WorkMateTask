package io_bound_task

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type TaskStatus int

const (
	TaskStatusUnknown TaskStatus = iota
	TaskStatusActive
	TaskStatusPending
	TaskStatusSuccess
	TaskStatusError
)

func (t TaskStatus) String() string {
	switch t {

	case TaskStatusActive:
		return "Active"
	case TaskStatusPending:
		return "Pending"
	case TaskStatusSuccess:
		return "Success"
	case TaskStatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

type TaskType int

const (
	TaskTypeUnknown TaskType = iota
	TaskTypeIOBound
)

func (t TaskType) String() string {
	switch t {
	case TaskTypeIOBound:
		return "IOBound"
	default:
		return "Unknown"
	}
}
func StringToTaskType(s string) TaskType {
	switch s {
	case "IOBound":
		return TaskTypeIOBound
	default:
		return TaskTypeUnknown
	}
}

type Task struct {
	Id             uuid.UUID     `json:"id"`
	Type           TaskType      `json:"task_type"`
	Status         TaskStatus    `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	StartedAt      time.Time     `json:"started_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	ProcessingTime time.Duration `json:"processing_time"`
}

func (t *Task) Run(ctx context.Context) error {
	client := &http.Client{}
	for i := 1; i < 20; i++ {
		time.Sleep(15 * time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", "https://google.com", nil)
		if err != nil {
			return err
		}

		_, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("Error performing request: %v\n", err)
		}
	}

	return nil
}

func (t *Task) Stop(ctx context.Context) error {
	fmt.Println("Stopping task", t.Id)
	return nil
}

func (t *Task) GetID() uuid.UUID {
	return t.Id
}
