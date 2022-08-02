package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "team3-task/internal/controller/http/v1"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	repo "team3-task/internal/repository/kafka"
	"team3-task/internal/utils"
	"team3-task/pkg/logging"
)

type InUseCase struct {
	taskUseCase          TaskUseCase
	taskApproversUseCase TaskApproversUseCase
	userUseCase          UserUseCase
	kafkaClient          *repo.KafkaProducers
	log                  *logging.ZeroLogger
}

func NewInUseCase(ti TaskDBRepoInterface, ta TaskApproversDBRepoInterface, ui UserDBRepoInterface, kafkaClient *repo.KafkaProducers, log *logging.ZeroLogger) *InUseCase {
	taskUC := NewTaskUseCase(ti)
	taUC := NewTaskApproversUseCase(ta)
	userUC := NewUserUseCase(ui, log)
	return &InUseCase{*taskUC, *taUC, *userUC, kafkaClient, log}
}

var _ v1.TaskHandlerInterface = (*InUseCase)(nil)

func (inUC *InUseCase) CreateTaskHandle(ctx context.Context, data []byte, authorEmail string) (int, error) {
	receivedTask := entity.Task{}
	var taskID int // taskID is success
	err := json.Unmarshal(data, &receivedTask)
	if err != nil {
		err = fmt.Errorf("usecase.CreateTaskHandle unmarshal receivedTask error: %v", err)
		return taskID, err
	}

	// add author information
	err = utils.CheckEmail(authorEmail)
	if err != nil {
		return taskID, errors.Wrapf(err, "email is incorrect, please, check author information")
	}
	user, err := inUC.userUseCase.CheckAndReturnUserByEmail(ctx, authorEmail)
	if err != nil {
		return taskID, errors.Wrapf(err, "something went wrong, can't create author as user, please, try again")
	}
	receivedTask.Author = user

	// first - write users
	approverSlice := make([]entity.User, 0, len(receivedTask.Approvers))
	for _, v := range receivedTask.Approvers {
		err = utils.CheckEmail(v.Email)
		if err != nil {
			return taskID, errors.Wrapf(err, "email is incorrect, please, check approver information: %v", v.Email)
		}
		curApprover, err := inUC.userUseCase.CheckAndReturnUserByEmail(ctx, v.Email)
		if err != nil {
			return taskID, errors.Wrapf(err, "something went wrong, can't create approver %v as user, please, try again", v.Email)
		}
		approverSlice = append(approverSlice, curApprover)
	}
	receivedTask.Approvers = approverSlice

	//-----------------------------------
	// begin transaction
	tx, err := inUC.taskUseCase.dbRepo.BeginTransaction(ctx)
	if err != nil {
		return taskID, errors.Wrapf(err, "something went wrong, can't create task, please, try again")
	}
	defer inUC.taskUseCase.dbRepo.RollbackTransaction(ctx, tx)

	// second - write task
	taskID, err = inUC.taskUseCase.CreateTask(ctx, receivedTask)
	if err != nil {
		inUC.taskUseCase.dbRepo.RollbackTransaction(ctx, tx)
		err = fmt.Errorf("in usecase.CreateTaskHandle inUC.taskUseCase.CreateTask error: %v", err)
		return taskID, err
	}

	// third - write task approvers to db
	err = inUC.taskApproversUseCase.InsertTaskApprovers(ctx, tx, taskID, receivedTask.Approvers)
	if err != nil {
		inUC.taskUseCase.dbRepo.RollbackTransaction(ctx, tx)
		err = fmt.Errorf("in usecase.CreateTaskHandle inUC.taskApproversUseCase.InsertTaskApprovers error: %v", err)
		return taskID, err
	}

	// fourth - write event to db (created)

	// send kafka message to analytic service about task create
	receivedTask.Id = taskID
	// последние два параметра - требования бизнес-логики сервиса аналитики - тип события (создание, согласование и т.д.) и пользователь (если создание, то email автора, если согласование, то email согласующего)
	err = repo.SendMessagesToKafka(inUC.kafkaClient.KafProducerAboutTaskEvent, &receivedTask, entity.Created, receivedTask.Author.Email, entity.AboutTaskEvent)
	if err != nil {
		inUC.taskUseCase.dbRepo.RollbackTransaction(ctx, tx)
		err = fmt.Errorf("usecase.CreateHandle repo.SendMessagesToKafka error: %v", err)
		return taskID, err
	}

	// commit transaction
	_ = inUC.taskUseCase.dbRepo.CommitTransaction(ctx, tx)
	//-----------------------------------

	return taskID, nil
}

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
