package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"
)

type TaskEventsUseCase struct {
	dbRepo TaskEventsDBRepoInterface
}

func NewTaskEventsUseCase(r TaskEventsDBRepoInterface) *TaskEventsUseCase {
	return &TaskEventsUseCase{r}
}

func (teUC *TaskEventsUseCase) InsertTaskEvent(ctx context.Context, taskID, userId int, taskEventType entity.KafkaTypes) error {

	err := teUC.dbRepo.InsertDBTaskEvents(ctx, taskID, userId, taskEventType)
	if err != nil {
		err = fmt.Errorf("usecase: (teUC *TaskEventsUseCase) InsertTaskEvent: InsertDBTaskEvents error: %v", err)
		return err
	}
	return nil
}
