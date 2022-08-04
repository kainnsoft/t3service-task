package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"

	"github.com/jackc/pgx/v4"
)

type TaskUseCase struct {
	dbRepo TaskDBRepoInterface
}

func NewTaskUseCase(r TaskDBRepoInterface) *TaskUseCase {
	return &TaskUseCase{r}
}

func (taskUC *TaskUseCase) CreateTask(ctx context.Context, txPtr *pgx.Tx, task *entity.Task) (int, error) {

	taskID, err := taskUC.dbRepo.CreateDBTask(ctx, txPtr, task)
	if err != nil {
		err = fmt.Errorf("usecase.CreateTask taskUC.dbRepo.CreateDBTask error: %v", err)
		return taskID, err
	}
	return taskID, nil
} // TODO
