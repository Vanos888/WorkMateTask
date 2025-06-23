package main

import (
	ogen_server "WorkMateTask/internal/generated/servers/http/v1/task"
	task2 "WorkMateTask/internal/handler/api/http/v1/task"
	task_processing "WorkMateTask/internal/jobs/task-processing"
	task_repository "WorkMateTask/internal/repository/task-repository"
	processing_service "WorkMateTask/internal/service/processing-service"
	task_service "WorkMateTask/internal/service/task-service"
	"WorkMateTask/internal/storage"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	stor := storage.NewStorage()

	repository := task_repository.NewTaskRepository(stor)

	processing, err := task_processing.NewTaskProcessor(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	taskProcessingService := processing_service.NewTaskProcessingService(processing, stor)
	service := task_service.NewTaskService(repository, taskProcessingService)

	api := task2.NewTaskApi(service)

	srv, err := ogen_server.NewServer(api)
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	if err = taskProcessingService.ObserveTasks(ctx); err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
