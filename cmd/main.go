package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fashionparty/fast-task-manager-server/internal/handler"
	"github.com/fashionparty/fast-task-manager-server/internal/middleware"
	"github.com/fashionparty/fast-task-manager-server/internal/model"
	"github.com/fashionparty/fast-task-manager-server/internal/repository"
	"github.com/fashionparty/fast-task-manager-server/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	databaseDSN  = "host=localhost port=5432 user=fasttaskmanager password=password dbname=fasttaskmanager_db sslmode=disable"
	secret       = "super-secret-key-aaaaaaaaaaaaaaa"
	expiration   = 24 * time.Hour
	serverURL    = ":8080"
	seedEmail    = "test@test.com"
	seedPassword = "testtest"
)

func main() {
	db, err := setupDatabase()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	authService := service.NewAuthService(userRepo, secret, expiration)
	taskService := service.NewTaskService(taskRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)

	authMiddleware := middleware.NewAuthMiddleware(secret)

	if err := seedDatabase(context.Background(), authService, taskService); err != nil {
		slog.Error("failed to seed database", "error", err)
		os.Exit(1)
	}

	mux := setupRouter(authHandler, taskHandler, authMiddleware)

	slog.Info("Server running", "address", serverURL)
	err = http.ListenAndServe(serverURL, mux)
	if err != nil {
		slog.Info("Server shutdown")
		os.Exit(1)
	}
}

func setupRouter(
	authHandler handler.AuthHandler,
	taskHandler handler.TaskHandler,
	authMiddleware *middleware.AuthMiddleware,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.Handle("POST /api/tasks", authMiddleware.Authenticate(http.HandlerFunc(taskHandler.CreateTask)))
	mux.Handle("GET /api/tasks", authMiddleware.Authenticate(http.HandlerFunc(taskHandler.ListTasks)))
	mux.Handle("GET /api/tasks/{id}", authMiddleware.Authenticate(http.HandlerFunc(taskHandler.GetTask)))
	mux.Handle("PUT /api/tasks/{id}", authMiddleware.Authenticate(http.HandlerFunc(taskHandler.UpdateTask)))
	mux.Handle("DELETE /api/tasks/{id}", authMiddleware.Authenticate(http.HandlerFunc(taskHandler.DeleteTask)))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			slog.Error("failed to write health check response", "error", err)
		}
	})

	return mux
}

func setupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	slog.Info("database connected successfully")

	if err := db.AutoMigrate(&model.User{}, &model.Task{}); err != nil {
		return nil, err
	}

	return db, nil
}

func seedDatabase(ctx context.Context, authService service.AuthService, taskService service.TaskService) error {
	user, err := authService.Register(ctx, seedEmail, seedPassword)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			return nil
		}
		return err
	}

	dueDate := time.Now().Add(48 * time.Hour)

	seedTasks := []struct {
		title    string
		content  string
		priority model.Priority
		dueDate  *time.Time
		isDone   bool
	}{
		{"Prepare presentation", "Slides for the Monday review", model.PriorityHigh, &dueDate, false},
		{"Buy groceries", "Milk, bread, coffee", model.PriorityLow, nil, false},
		{"Answer emails", "Backlog from Friday", model.PriorityMedium, nil, true},
	}

	for _, seedTask := range seedTasks {
		task, err := taskService.CreateTask(ctx, user.ID, seedTask.title, seedTask.content, seedTask.priority, seedTask.dueDate)
		if err != nil {
			return err
		}

		if seedTask.isDone {
			if _, err := taskService.UpdateTask(ctx, user.ID, task.ID, nil, nil, nil, &seedTask.isDone, nil); err != nil {
				return err
			}
		}
	}

	slog.Info("seed user created", "email", seedEmail, "tasks", len(seedTasks))

	return nil
}

func getDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return databaseDSN
}
