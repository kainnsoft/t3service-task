package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"

	"github.com/jackc/pgx/v4"
)

type TaskUseCase struct {
	dbRepo   TaskDBRepoInterface
	teDBRepo TaskEventsDBRepoInterface
}

func NewTaskUseCase(r TaskDBRepoInterface, te TaskEventsDBRepoInterface) *TaskUseCase {
	return &TaskUseCase{r, te}
}

func (taskUC *TaskUseCase) CreateTask(ctx context.Context, txPtr *pgx.Tx, task *entity.Task) (int, error) {
	taskID, err := taskUC.dbRepo.CreateDBTask(ctx, txPtr, task)
	if err != nil {
		err = fmt.Errorf("usecase.CreateTask taskUC.dbRepo.CreateDBTask error: %w", err)

		return taskID, err
	}

	return taskID, nil
}

func (taskUC *TaskUseCase) InsertTaskEvent(ctx context.Context, taskID, userID int, taskEventType entity.KafkaTypes) error {
	err := taskUC.teDBRepo.InsertDBTaskEvents(ctx, taskID, userID, taskEventType)
	if err != nil {
		err = fmt.Errorf("usecase.InsertTaskEvent: taskUC.teDBRepo.InsertDBTaskEvents error: %w", err)

		return err
	}

	return nil
}
