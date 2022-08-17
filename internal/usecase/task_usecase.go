package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"
	"team3-task/pkg/logging"

	"github.com/jackc/pgx/v4"
)

type TaskUseCase struct {
	dbRepo     TaskDBRepoInterface
	userDBRepo UserDBRepoInterface
	taDBRepo   TaskApproversDBRepoInterface
	teDBRepo   TaskEventsDBRepoInterface
	log        *logging.ZeroLogger
}

func NewTaskUseCase(r TaskDBRepoInterface,
	usRepo UserDBRepoInterface,
	ta TaskApproversDBRepoInterface,
	te TaskEventsDBRepoInterface,
	log *logging.ZeroLogger) *TaskUseCase {

	carTaskUseCase := TaskUseCase{r, usRepo, ta, te, log}
	return &carTaskUseCase
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

func (taskUC *TaskUseCase) GetTaskList(ctx context.Context) ([]entity.Task, error) {
	taskList, err := taskUC.dbRepo.GetListDBTask(ctx) // without approvers
	if err != nil {
		err = fmt.Errorf("usecase.GetTaskList: taskUC.dbRepo.GetListDBTask error: %w", err)

		return nil, err
	}

	for i, t := range taskList {
		// set authors
		author, err := taskUC.userDBRepo.GetDBUserByID(ctx, t.Author.ID)
		if err != nil {
			taskUC.log.Error("usecase.GetTaskList (taskID = %d): taskUC.userDBRepo.GetDBUserByID error :%w", t.ID, err)
		}

		taskList[i].Author = author

		// add approvers
		apprList, err := taskUC.taDBRepo.GetTaskApproversByTaskID(ctx, t.ID)
		if err != nil {
			taskUC.log.Error("usecase.GetTaskList (taskID = %d): taskUC.taDBRepo.GetTaskApproversByTaskID error :%w", t.ID, err)
			continue
		}

		taskList[i].Approvers = apprList
	}

	return taskList, nil
}

func (taskUC *TaskUseCase) GetOneTask(ctx context.Context, taskID int) (entity.Task, error) {
	task := entity.Task{}
	task, err := taskUC.dbRepo.GetDBTask(ctx, taskID) // without approvers
	if err != nil {
		err = fmt.Errorf("usecase.GetOneTask: taskUC.dbRepo.GetDBTask error: %w", err)

		return task, err
	}

	// set authors
	author, err := taskUC.userDBRepo.GetDBUserByID(ctx, task.Author.ID)
	if err != nil {
		taskUC.log.Error("usecase.GetOneTask (taskID = %d): taskUC.userDBRepo.GetDBUserByID error :%w", task.ID, err)
	}

	task.Author = author

	// add approvers
	apprList, err := taskUC.taDBRepo.GetTaskApproversByTaskID(ctx, task.ID)
	if err != nil {
		taskUC.log.Error("usecase.GetOneTask (taskID = %d): taskUC.taDBRepo.GetTaskApproversByTaskID error :%w", task.ID, err)
	}

	task.Approvers = apprList

	return task, nil
}
