package service

import (
	"context"
	"time"

	"github.com/fashionparty/fast-task-manager-server/internal/model"
	"github.com/fashionparty/fast-task-manager-server/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, userID uint, title, content string, priority model.Priority, dueDate *time.Time) (*model.Task, error)
	GetTask(ctx context.Context, userID, taskID uint) (*model.Task, error)
	UpdateTask(ctx context.Context, userID, taskID uint, title, content *string, priority *model.Priority, isDone *bool, dueDate *time.Time) (*model.Task, error)
	ListTasks(ctx context.Context, userID uint) ([]model.Task, error)
	DeleteTask(ctx context.Context, userID, taskID uint) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateTask(ctx context.Context, userID uint, title, content string, priority model.Priority, dueDate *time.Time) (*model.Task, error) {
	newTask := &model.Task{
		Title:    title,
		Content:  content,
		DueDate:  dueDate,
		Priority: priority,
		IsDone:   false,
		OwnerID:  userID,
	}

	err := s.taskRepo.Create(ctx, newTask)
	if err != nil {
		return nil, err
	}

	return newTask, nil
}

func (s *taskService) GetTask(ctx context.Context, userID, taskID uint) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.OwnerID != userID {
		return nil, repository.ErrNotFound
	}

	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, userID, taskID uint, title, content *string, priority *model.Priority, isDone *bool, dueDate *time.Time) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.OwnerID != userID {
		return nil, repository.ErrNotFound
	}

	if title != nil {
		task.Title = *title
	}
	if content != nil {
		task.Content = *content
	}
	if priority != nil {
		task.Priority = *priority
	}
	if isDone != nil {
		task.IsDone = *isDone
	}
	if dueDate != nil {
		task.DueDate = dueDate
	}

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, userID, taskID uint) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.OwnerID != userID {
		return repository.ErrNotFound
	}

	return s.taskRepo.Delete(ctx, taskID)
}

func (s *taskService) ListTasks(ctx context.Context, userID uint) ([]model.Task, error) {
	return s.taskRepo.List(ctx, userID)
}
