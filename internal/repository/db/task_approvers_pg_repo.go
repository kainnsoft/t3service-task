package repository

import (
	"context"
	"team3-task/internal/entity"
	"team3-task/internal/usecase"
	"team3-task/pkg/pg"

	"github.com/jackc/pgx/v4"
)

type TaskApproversPGRepo struct {
	*pg.DB
}

var _ usecase.TaskApproversDBRepoInterface = (*TaskApproversPGRepo)(nil)

func NewTaskApproversPGRepo(pgdb *pg.DB) *TaskApproversPGRepo {
	return &TaskApproversPGRepo{pgdb}
}

func (repo *TaskApproversPGRepo) InsertDBTaskApprovers(ctx context.Context, txPtr *pgx.Tx, taskId int, userList []entity.User) error {
	const query = "INSERT INTO task.task_approvers (task_id, approver_id) VALUES ($1, $2);"
	batch := new(pgx.Batch)

	for _, user := range userList {
		batch.Queue(query, taskId, user.ID)
	}

	tx := *txPtr
	res := tx.SendBatch(ctx, batch)

	err := res.Close()
	if err != nil {
		return err
	}

	return nil
}

func (repo *TaskApproversPGRepo) GetTaskApproversByTaskID(ctx context.Context, taskID int) ([]entity.User, error) {
	userList := make([]entity.User, 0)

	const queryStrID = "select approver_id from task_approvers where task_id=$1;"
	rows, err := repo.Pool.Query(ctx, queryStrID, taskID)
	if err == pgx.ErrNoRows {
		return userList, err
	} else if err != nil {
		return userList, err
	}
	defer rows.Close()

	userIDSlice := make([]int, 0)

	for rows.Next() {
		var curID int
		err = rows.Scan(&curID)
		if err != nil {
			return userList, err
		}

		userIDSlice = append(userIDSlice, curID)
	}

	const queryStrUser = "select id, email from users where id = ANY($1::int[]) ORDER BY id ASC;"
	// param := "{" + strings.Join(userIDSlice, ",") + "}"
	rowsN, err := repo.Pool.Query(ctx, queryStrUser, userIDSlice)
	if err == pgx.ErrNoRows {
		return userList, err
	} else if err != nil {
		return userList, err
	}
	defer rowsN.Close()

	for rows.Next() {
		curUser := entity.User{}
		err = rows.Scan(&curUser)
		if err != nil {
			return userList, err
		}

		userList = append(userList, curUser)
	}

	return userList, nil
}

func (repo *TaskApproversPGRepo) GetTaskApproversIDByTaskID(ctx context.Context, taskID int) ([]int, error) {
	approverIDList := make([]int, 0)

	const queryStrID = "select approver_id from task_approvers where task_id=$1;"
	rows, err := repo.Pool.Query(ctx, queryStrID, taskID)
	if err == pgx.ErrNoRows {
		return approverIDList, err
	} else if err != nil {
		return approverIDList, err
	}
	defer rows.Close()

	for rows.Next() {
		var curID int
		err = rows.Scan(&curID)
		if err != nil {
			return approverIDList, err
		}

		approverIDList = append(approverIDList, curID)
	}

	return approverIDList, nil
}
