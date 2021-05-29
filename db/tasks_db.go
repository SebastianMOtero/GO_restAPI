package db

import (
	"context"
	"strconv"

	"../models"
	"../rds"
	"github.com/jackc/pgx"
)

type TasksDB struct {
	dbClient rds.RdsClient
}

type TasksDBInterface interface {
	GetTasks(ctx context.Context, id int) ([]models.Task, error)
	InsertTask(ctx context.Context, task models.Task) (*int, error)
	UpdateTask(ctx context.Context, task models.Task) (*int, error)
	DeleteTask(ctx context.Context, id int) (*int, error)
}

func NewTasksDB(dbClient rds.RdsClient) TasksDBInterface {
	return &TasksDB{
		dbClient: dbClient,
	}
}

func (db TasksDB) GetTasks(ctx context.Context, id int) ([]models.Task, error) {
	var taskList []models.Task
	var conn rds.Querier
	var err error
	var rows pgx.Rows

	if conn, err = db.dbClient.Querier(ctx, nil); err != nil {
		return nil, err
	}

	sql := `SELECT * FROM tasks`
	if id != 0 {
		sql = sql + ` WHERE id=` + strconv.Itoa(id)
	}
	if rows, err = conn.Query(ctx, sql); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		if err := rows.Scan(&task.ID, &task.Name, &task.Content); err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func (db TasksDB) InsertTask(ctx context.Context, task models.Task) (*int, error) {
	var conn rds.Querier
	var err error
	var id int

	if conn, err = db.dbClient.Querier(ctx, nil); err != nil {
		return nil, err
	}

	sql := `INSERT INTO tasks(name, content) VALUES($1, $2) RETURNING id`

	err = conn.QueryRow(ctx, sql, task.Name, task.Content).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (db TasksDB) UpdateTask(ctx context.Context, task models.Task) (*int, error) {
	var conn rds.Querier
	var err error
	var id int

	if conn, err = db.dbClient.Querier(ctx, nil); err != nil {
		return nil, err
	}

	sql := `UPDATE tasks SET name= $1, content=$2 WHERE id=$3 RETURNING id`

	err = conn.QueryRow(ctx, sql, task.Name, task.Content, task.ID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (db TasksDB) DeleteTask(ctx context.Context, id int) (*int, error) {
	var conn rds.Querier
	var err error

	if conn, err = db.dbClient.Querier(ctx, nil); err != nil {
		return nil, err
	}

	sql := `DELETE FROM tasks WHERE id=$1 RETURNING id`

	err = conn.QueryRow(ctx, sql, id).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
