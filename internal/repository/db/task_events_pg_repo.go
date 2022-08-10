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

func NewTaskEventsPGRepo(pgdb *pg.DB) *TaskEventsPGRepo {
	return &TaskEventsPGRepo{pgdb}
}

func (repo *TaskEventsPGRepo) InsertDBTaskEvents(ctx context.Context, taskID, userID int, taskEventType entity.KafkaTypes) error {
	queryStr := "INSERT INTO task.task_events (task_id, user_id, event_type_id, event_time) VALUES ($1, $2, $3, $4);"

	taskEventTypeID, err := repo.GetTaskEventTypeByName(ctx, taskEventType)
	if err != nil {
		return fmt.Errorf("repository (repo *TaskEventsPGRepo) GetTaskEventTypeByName error: %v", err)
	}

	_, err = repo.Pool.Exec(ctx, queryStr,
		taskID,
		userID,
		taskEventTypeID,
		time.Now())
	if err != nil {
		return errors.Wrapf(err, "repository (repo *TaskEventsPGRepo) InsertDBTaskEvents error")
	}

	return nil
}

func (repo *TaskEventsPGRepo) GetApproversIDMapMatchingTheListByTaskID(ctx context.Context, taskID int, listApproversID []int) (map[int]struct{}, error) {
	eventApproversIDMap := make(map[int]struct{})
	queryStr := "select id, user_id from task_events where task_id = ANY($1::int[]);"
	// param := "{" + strings.Join(userIDSlice, ",") + "}"
	rows, err := repo.Pool.Query(ctx, queryStr, listApproversID)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var curApproverID int
		rows.Scan(&curApproverID)
		eventApproversIDMap[curApproverID] = struct{}{}
	}

	return eventApproversIDMap, nil
}

func (repo *TaskEventsPGRepo) GetTaskEventTypeByName(ctx context.Context, taskEventType entity.KafkaTypes) (int, error) {
	const queryStr = "select id from task.task_event_types where task_type = $1"
	row := repo.Pool.QueryRow(ctx, queryStr, taskEventType)

	var taskTypeID int
	err := row.Scan(&taskTypeID)
	if err == pgx.ErrNoRows {
		return taskTypeID, err
	} else if err != nil {
		return taskTypeID, err
	}

	return taskTypeID, nil
}
