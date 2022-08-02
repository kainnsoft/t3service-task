package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"

	"github.com/jackc/pgx/v4"
)

type TaskApproversUseCase struct {
	dbRepo TaskApproversDBRepoInterface
}

func NewTaskApproversUseCase(r TaskApproversDBRepoInterface) *TaskApproversUseCase {
	return &TaskApproversUseCase{r}
}

func (taUC *TaskApproversUseCase) InsertTaskApprovers(ctx context.Context, txPtr *pgx.Tx, taskId int, approverList []entity.User) error {

	err := taUC.dbRepo.InsertDBTaskApprovers(ctx, txPtr, taskId, approverList)
	if err != nil {
		err = fmt.Errorf("usecase.InsertTaskApprovers taUC.dbRepo.InsertDBTaskApprovers error: %v", err)
		return err
	}
	return nil
}
