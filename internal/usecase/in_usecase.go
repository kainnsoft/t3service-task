package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "team3-task/internal/controller/http/v1"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/utils"
	"team3-task/pkg/logging"
)

type InUseCase struct {
	taskUseCase TaskUseCase
	userUseCase UserUseCase
	log         *logging.ZeroLogger
}

func NewInUseCase(ti TaskDBRepoInterface, ui UserDBRepoInterface, log *logging.ZeroLogger) *InUseCase {
	taskUC := NewTaskUseCase(ti)
	userUC := NewUserUseCase(ui, log)
	return &InUseCase{*taskUC, *userUC, log}
}

var _ v1.TaskHandlerInterface = (*InUseCase)(nil)

func (inUC *InUseCase) CreateTaskHandle(ctx context.Context, data []byte, authorEmail string) (string, error) {
	receivedTask := entity.Task{}
	var resp string // taskID is success
	err := json.Unmarshal(data, &receivedTask)
	if err != nil {
		err = fmt.Errorf("usecase.CreateHandle unmarshal receivedTask error: %v", err)
		return resp, err
	}

	// add author information
	err = utils.CheckEmail(authorEmail)
	if err != nil {
		return "", errors.Wrapf(err, "email is incorrect, please, check author information")
	}
	user, err := inUC.userUseCase.CheckAndReturnUserByEmail(ctx, authorEmail)
	if err != nil {
		return "", errors.Wrapf(err, "something went wrong, can't create author as user, please, try again")
	}
	receivedTask.Author = user
	//-----------------------------------
	// begin transaction
	// first - write users
	approverSlice := make([]entity.User, 0, len(receivedTask.Approvers))
	for _, v := range receivedTask.Approvers {
		err = utils.CheckEmail(v.Email)
		if err != nil {
			return "", errors.Wrapf(err, "email is incorrect, please, check approver information: %v", v.Email)
		}
		curApprover, err := inUC.userUseCase.CheckAndReturnUserByEmail(ctx, v.Email)
		if err != nil {
			return "", errors.Wrapf(err, "something went wrong, can't create approver %v as user, please, try again", v.Email)
		}
		approverSlice = append(approverSlice, curApprover)
	}
	receivedTask.Approvers = approverSlice
	// second - write task
	resp, err = inUC.taskUseCase.CreateTask(ctx, receivedTask)
	if err != nil {
		err = fmt.Errorf("usecase.CreateHandle taskUC.dbRepo.Create receivedTask error: %v", err)
		return resp, err
	}
	// third - write task approvers
	// fourth - write event (created)
	// commit transaction
	//-----------------------------------
	return resp, nil
} // TODO

func (taskUC *InUseCase) UpdateTaskHandle(ctx context.Context, task entity.Task) (int, error) {
	// resp, err := taskUC.dbRepo.UpdateDBTask(context.Background(), task)
	// if err != nil {
	// 	return 0, err
	// }
	return 0, nil
} // TODO

func (taskUC *InUseCase) DeleteTaskHandle(ctx context.Context, taskId int) error {
	// err := taskUC.dbRepo.DeleteDBTask(context.Background(), taskId)
	// if err != nil {
	// 	return err
	// }
	return nil
} // TODO

func (taskUC *InUseCase) GetTaskHandle(ctx context.Context, taskId int) (entity.Task, error) { // TODO
	// resp, err := taskUC.dbRepo.GetDBTask(context.Background(), taskId)
	// if err != nil {
	// 	return entity.Task{}, err
	// }
	return entity.Task{}, nil
} // TODO

func (taskUC *InUseCase) ListTaskHandle(ctx context.Context) ([]entity.Task, error) {
	emptytaskList := make([]entity.Task, 0)
	// taskList, err := taskUC.dbRepo.ListDBTask(context.Background())
	// if err != nil {
	// 	return emptytaskList, err
	// }
	return emptytaskList, nil
} // TODO (with filtr)
