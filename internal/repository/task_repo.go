package repository

import (
	"database/sql"
	"todo-app/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (t *TaskRepository) Create(input models.CreateTaskInput) (models.Task, error) {

	var task models.Task

	sqlQuery := `
	INSERT INTO tasks (title, description)
	VALUES ($1, $2)
	RETURNING id, title, description, is_completed, created_at, completed_at;
	`

	err := t.db.QueryRow(sqlQuery, input.Title, input.Description).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.IsCompleted,
		&task.CreatedAt,
		&task.CompletedAt,
	)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil

}

func (t *TaskRepository) GetAll() ([]models.Task, error) {
	outputTasks := []models.Task{}

	sqlQuery := `
	SELECT id, title, description, is_completed, created_at, completed_at
	FROM tasks 
	ORDER BY id;
	`

	rows, err := t.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsCompleted,
			&task.CreatedAt,
			&task.CompletedAt,
		); err != nil {
			return nil, err
		}
		outputTasks = append(outputTasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return outputTasks, nil

}

func (t *TaskRepository) GetByID(id int) (models.Task, error) {
	var task models.Task

	sqlQuery := `
	SELECT id, title, description, is_completed, created_at, completed_at
	FROM tasks 
	WHERE id=$1
	`

	if err := t.db.QueryRow(sqlQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.IsCompleted,
		&task.CreatedAt,
		&task.CompletedAt,
	); err != nil {
		return models.Task{}, err
	}

	return task, nil

}

func (t *TaskRepository) Delete(id int) error {
	sqlQuery := `
	DELETE FROM tasks
	WHERE id=$1
	`

	if _, err := t.db.Exec(sqlQuery, id); err != nil {
		return err
	}

	return nil
}

func (t *TaskRepository) Update(inputTask models.UpdateTaskInput) (models.Task, error) {
	var task models.Task

	sqlQuery := `
	UPDATE tasks 
    SET title = $1, description = $2, is_completed = $3, 
    completed_at = CASE WHEN $3 = true THEN NOW() ELSE NULL END
    WHERE id = $4 
    RETURNING id, title, description, is_completed, created_at, completed_at
	`

	if err := t.db.QueryRow(sqlQuery, inputTask.Title, inputTask.Description, inputTask.IsCompleted, inputTask.ID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.IsCompleted,
		&task.CreatedAt,
		&task.CompletedAt,
	); err != nil {
		return models.Task{}, err
	}

	return task, nil
}
