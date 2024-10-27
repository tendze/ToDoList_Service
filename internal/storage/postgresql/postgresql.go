package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
	"todolist/internal/storage"
)

type Storage struct {
	DB *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.New"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	DO $$ BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_status') THEN
			CREATE TYPE task_status AS ENUM('new', 'done');
		END IF;
	END $$;
`)
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    user_login VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
	status task_status DEFAULT 'new',
    created_at TIMESTAMP DEFAULT NOW(),
    deadline TIMESTAMP DEFAULT NOW()
)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	strg := &Storage{DB: db}
	return strg, nil
}

func (s *Storage) AddTask(userLogin, title, description string, deadline time.Time) error {
	const op = "storage.postgresql.AddTask"
	query := `INSERT INTO tasks(user_login, title, description, deadline) VALUES($1, $2, $3, $4)`
	_, err := s.DB.Exec(
		query,
		userLogin,
		title,
		description,
		deadline,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CompleteTask(id int, login string, newTaskStatus storage.TaskStatus) error {
	const op = "storage.postgresql.CompleteTask"
	query := `UPDATE tasks SET status = $1 WHERE id = $2 AND user_login = $3`
	_, err := s.DB.Exec(query, newTaskStatus, id, login)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
