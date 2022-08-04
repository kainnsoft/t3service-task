package repository

import (
	"context"
	"fmt"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/usecase"
	"team3-task/pkg/pg"
	"time"

	"github.com/jackc/pgx/v4"
)

type TaskEventsPGRepo struct {
	*pg.DB
}

var _ usecase.TaskEventsDBRepoInterface = (*TaskEventsPGRepo)(nil)

func NewTaskEventsPGRepo(pg *pg.DB) *TaskEventsPGRepo {
	return &TaskEventsPGRepo{pg}
}

func (repo *TaskEventsPGRepo) InsertDBTaskEvents(ctx context.Context, taskID int, userId int, taskEventType entity.KafkaTypes) error {
	queryStr := "INSERT INTO task.task_events (task_id, user_id, event_type_id, event_time) VALUES ($1, $2, $3, $4);"

	taskEventTypeId, err := repo.GetTaskEventTypeByName(ctx, taskEventType)
	if err != nil {
		return fmt.Errorf("repository (repo *TaskEventsPGRepo) GetTaskEventTypeByName error: %v", err)
	}

	_, err = repo.Pool.Exec(ctx, queryStr,
		taskID,
		userId,
		taskEventTypeId,
		time.Now())
	if err != nil {
		return errors.Wrapf(err, "repository (repo *TaskEventsPGRepo) InsertDBTaskEvents error")
	}
	return nil
}

func (repo *TaskEventsPGRepo) GetTaskEventTypeByName(ctx context.Context, taskEventType entity.KafkaTypes) (int, error) {
	const queryStr = "select id from task.task_event_types where task_type = $1"
	row := repo.Pool.QueryRow(ctx, queryStr, taskEventType)

	var taskTypeId int
	err := row.Scan(&taskTypeId)
	if err == pgx.ErrNoRows {
		return taskTypeId, err
	} else if err != nil {
		return taskTypeId, err
	}
	return taskTypeId, nil
}
