package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	config "todo-app/config"
	"todo-app/internal/handlers"
	"todo-app/internal/repository"

	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg := config.LoadConfig()
	fmt.Println("DSN:", cfg.DBURL)
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		sugar.Fatalw("Ошибка конфигурации DSN", "error", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		sugar.Fatalw("Не удалось подключиться к базе:", "error", err)
	}

	sugar.Info("Успешное подключение к PostgreSQL!")

	repo := repository.NewTaskRepository(db)
	handler := &handlers.TaskHandler{TaskRepository: repo, Logger: sugar}

	r := mux.NewRouter()
	r.HandleFunc("/tasks", handler.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", handler.GetAllHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.UpdateHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", handler.DeleteHandler).Methods("DELETE")

	sugar.Info("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
