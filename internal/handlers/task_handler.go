package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"todo-app/internal/models"
	"todo-app/internal/repository"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type TaskHandler struct {
	TaskRepository *repository.TaskRepository
	Logger         *zap.SugaredLogger
}

func (t *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	task := models.CreateTaskInput{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("некорректное тело запроса", "error", err)
		w.Write([]byte(err.Error()))
		return
	}
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("попытка создать задачу с пустым названием")
		w.Write([]byte("title is required!!!"))
		return
	}
	ret, err := t.TaskRepository.Create(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.Logger.Errorw("Ошибка при создании задачи", "error", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		t.Logger.Errorw("ошибка отправки ответа", "error", err)
		return
	}
	t.Logger.Infow("задача создана успешно.",
		"id", ret.ID, "title", ret.Title, "description", ret.Description)

}

func (t *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	str, exist := vars["id"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("ошибка. невозможно извлечь параметр id")
		return
	}
	idInt, err := strconv.Atoi(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("ошибка. id должен иметь числовой тип int", "error", err)
		w.Write([]byte(err.Error()))
		return
	}

	outputTask, err := t.TaskRepository.GetByID(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			t.Logger.Warnw("ошибка. пользователя с таким id не существует", "error", err)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		t.Logger.Warnw("ошибка. невозможно получить пользователя с таким id")
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(outputTask); err != nil {
		t.Logger.Errorw("ошибка отправки ответа", "error", err)
		return
	}
}

func (t *TaskHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	model := models.UpdateTaskInput{}
	vars := mux.Vars(r)
	str, exist := vars["id"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("ошибка. невозможно извлечь параметр id")
		return
	}
	idInt, err := strconv.Atoi(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("ошибка. id должен иметь числовой тип int", "error", err)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("некорректное тело запроса", "error", err)
		w.Write([]byte(err.Error()))
		return
	}
	model.ID = idInt

	outputModel, err := t.TaskRepository.Update(model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			t.Logger.Warnw("ошибка. задачи, которую нужно обновить, нет", "error", err)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		t.Logger.Errorw("ошибка обнвления задачи")
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(outputModel); err != nil {
		t.Logger.Errorw("ошибка отправки ответа", "error", err)
		return
	}
	t.Logger.Infow("задача обновлена успешно.", "id:", outputModel.ID,
		"title", outputModel.Title, "description", outputModel.Description, "status", outputModel.IsCompleted)
}

func (t *TaskHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	str, exist := vars["id"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("ошибка. невозможно извлечь параметр id")
		return
	}
	idInt, err := strconv.Atoi(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		t.Logger.Warnw("ошибка. id должен иметь числовой тип int", "error", err)
		w.Write([]byte(err.Error()))
		return
	}

	if err := t.TaskRepository.Delete(idInt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			t.Logger.Warnw("ошибка. нет задачи с таким id для удаления", "error", err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		t.Logger.Errorw("ошибка удаления задачи", "error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	t.Logger.Infow("задача удалена успешно.", "id", idInt)

}

func (t *TaskHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	outputModel, err := t.TaskRepository.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		t.Logger.Warnw("проблема с запросом", "error", err)
		return
	}
	if err := json.NewEncoder(w).Encode(outputModel); err != nil {
		t.Logger.Errorw("ошибка отправки ответа", "error", err)
		return
	}
}
