package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"
)

type TaskUseCase struct {
	dbRepo TaskDBRepoInterface
}

func NewTaskUseCase(r TaskDBRepoInterface) *TaskUseCase {
	return &TaskUseCase{r}
}

func (taskUC *TaskUseCase) CreateTask(ctx context.Context, task entity.Task) (string, error) {

	resp, err := taskUC.dbRepo.CreateDBTask(context.Background(), task)
	if err != nil {
		err = fmt.Errorf("usecase.CreateHandle taskUC.dbRepo.Create receivedTask error: %v", err)
		return resp, err
	}
	return resp, nil
} // TODO
