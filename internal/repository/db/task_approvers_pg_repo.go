package repository

import (
	"context"
	"team3-task/internal/entity"
	"team3-task/internal/usecase"
	"team3-task/pkg/logging"
	"team3-task/pkg/pg"

	"github.com/jackc/pgx/v4"
)

type TaskApproversPGRepo struct {
	logger *logging.ZeroLogger
	*pg.PgDB
}

var _ usecase.TaskApproversDBRepoInterface = (*TaskApproversPGRepo)(nil)

func NewTaskApproversPGRepo(pg *pg.PgDB, logger *logging.ZeroLogger) *TaskApproversPGRepo {
	return &TaskApproversPGRepo{logger, pg}
}

func (repo *TaskApproversPGRepo) InsertDBTaskApprovers(ctx context.Context, txPtr *pgx.Tx, taskId int, userList []entity.User) error {
	query := "INSERT INTO task.task_approvers (task_id, approver_id) VALUES ($1, $2);"
	batch := new(pgx.Batch)

	for _, user := range userList {
		batch.Queue(query, taskId, user.Id)
	}

	tx := *txPtr
	res := tx.SendBatch(ctx, batch)

	err := res.Close()
	if err != nil {
		return err
	}
	return nil
}

func (repo *TaskApproversPGRepo) GetTaskApproversByTaskID(ctx context.Context, taskId int) ([]entity.User, error) {
	userList := make([]entity.User, 0)
	//repo.Pool.Exec()
	return userList, nil // TODO
}
