// методы работы с БД. Реализация интерфейсов, которые будут "дергаться" в usecase-ах
package repository

import (
	"context"
	"strings"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/usecase"
	"team3-task/pkg/logging"
	"team3-task/pkg/pg"

	"github.com/jackc/pgx/v4"
)

type TaskPGRepo struct {
	logger *logging.ZeroLogger
	*pg.PgDB
}

var _ usecase.TaskDBRepoInterface = (*TaskPGRepo)(nil)

func NewTaskPGRepo(pg *pg.PgDB, logger *logging.ZeroLogger) *TaskPGRepo {
	return &TaskPGRepo{logger, pg}
}

func (repo *TaskPGRepo) CreateDBTask(ctx context.Context, task entity.Task) (int, error) {

	var taskID int
	query := "INSERT INTO task.tasks (author_id, descr, body, finished) VALUES ($1, $2, $3, $4) RETURNING id;"

	err := repo.Pool.QueryRow(ctx, query,
		task.Author.Id,
		strings.TrimSpace(task.Descr),
		strings.TrimSpace(task.Body),
		false).Scan(&taskID)
	if err != nil {
		return taskID, errors.Wrapf(err, "repository (repo *TaskDBRepo) Create error") // AddErrorContext(err, "task", " task should't be empty") //.
	}
	return taskID, nil
}

func (repo *TaskPGRepo) UpdateDBTask(ctx context.Context, task entity.Task) (int, error) {
	//repo.Pool.Exec()
	return 0, nil // TODO
}

func (repo *TaskPGRepo) DeleteDBTask(ctx context.Context, id int) error {
	//repo.Pool.Exec()
	return nil // TODO
}

func (repo *TaskPGRepo) GetDBTask(ctx context.Context, id int) (entity.Task, error) {

	const queryStr = "select id, author_id, descr, body, finished from task.tasks where id = $1"
	row := repo.Pool.QueryRow(ctx, queryStr, id)

	task := entity.Task{}
	err := row.Scan(&task.Id, &task.Author.Id)
	if err == pgx.ErrNoRows {
		return task, err
	} else if err != nil {
		return task, err
	}

	return entity.Task{}, nil
}

func (repo *TaskPGRepo) ListDBTask(ctx context.Context) ([]entity.Task, error) {
	//repo.Pool.Exec()
	return []entity.Task{}, nil // TODO
}

func (repo *TaskPGRepo) BeginTransaction(ctx context.Context) (*pgx.Tx, error) {
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (repo *TaskPGRepo) CommitTransaction(ctx context.Context, txPtr *pgx.Tx) error {
	tx := *txPtr
	return tx.Commit(ctx)
}

func (repo *TaskPGRepo) RollbackTransaction(ctx context.Context, txPtr *pgx.Tx) {
	tx := *txPtr
	_ = tx.Rollback(ctx)
}
