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
	taskEventsUseCase    TaskEventsUseCase
	userUseCase          UserUseCase
	txUseCase            TxUseCase
	kafkaClient          *repo.KafkaProducers
	log                  *logging.ZeroLogger
}

func NewInUseCase(ti TaskDBRepoInterface,
	ta TaskApproversDBRepoInterface,
	te TaskEventsDBRepoInterface,
	ui UserDBRepoInterface,
	tx TxDBRepoInterface,
	kafkaClient *repo.KafkaProducers,
	log *logging.ZeroLogger) *InUseCase {

	taskUC := NewTaskUseCase(ti)
	taUC := NewTaskApproversUseCase(ta)
	teUC := NewTaskEventsUseCase(te)
	userUC := NewUserUseCase(ui, log)
	txUC := NewTxUseCase(tx)

	return &InUseCase{*taskUC, *taUC, *teUC, *userUC, *txUC, kafkaClient, log}
}

var _ v1.TaskHandlerInterface = (*InUseCase)(nil)

func (inUC *InUseCase) CreateTaskHandle(ctx context.Context, data []byte, authorEmail string) (int, error) {
	var taskID int // taskID is success

	receivedTask := entity.Task{}

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
	tx, err := inUC.txUseCase.BeginTransaction(ctx)
	if err != nil {
		return taskID, errors.Wrapf(err, "something went wrong, can't create task, please, try again")
	}

	defer func() {
		err = inUC.txUseCase.RollbackTransaction(ctx, tx)
		if err != nil {
			inUC.log.Error("in usecase.CreateTaskHandle defer inUC.txUseCase.RollbackTransaction error :%v", err)
		}
	}()

	// second - write task
	taskID, err = inUC.taskUseCase.CreateTask(ctx, tx, &receivedTask)
	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle inUC.txUseCase.CreateTask RollbackTransaction error :%v", rberr)
		}
		err = fmt.Errorf("in usecase.CreateTaskHandle inUC.taskUseCase.CreateTask error: %v", err)

		return taskID, err
	}

	// third - write task approvers to db
	err = inUC.taskApproversUseCase.InsertTaskApprovers(ctx, tx, taskID, receivedTask.Approvers)
	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle inUC.txUseCase.InsertTaskApprovers RollbackTransaction error :%v", rberr)
		}
		err = fmt.Errorf("in usecase.CreateTaskHandle inUC.taskApproversUseCase.InsertTaskApprovers error: %v", err)

		return taskID, err
	}
	receivedTask.ID = taskID

	// fourth - write event to db (created)
	err = inUC.taskEventsUseCase.InsertTaskEvent(ctx, taskID, receivedTask.Author.ID, entity.Created)
	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle inUC.taskEventsUseCase.InsertTaskEvent RollbackTransaction error :%v", rberr)
		}
		err = fmt.Errorf("in usecase.CreateTaskHandle inUC.taskEventsUseCase.InsertTaskEvent error: %v", err)

		return taskID, err
	}

	// send kafka message to analytic service about task create
	// последние два параметра - требования бизнес-логики сервиса аналитики - тип события (создание, согласование и т.д.)
	// и пользователь (если создание, то email автора, если согласование, то email согласующего)
	err = repo.SendMessagesToKafka(
		inUC.kafkaClient.KafProducerAboutTaskEvent,
		&receivedTask,
		entity.Created,
		receivedTask.Author.Email,
		entity.AboutTaskEvent)

	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle repo.SendMessagesToKafka RollbackTransaction error :%v", rberr)
		}
		err = fmt.Errorf("usecase.CreateHandle repo.SendMessagesToKafka error: %v", err)

		return taskID, err
	}

	// commit transaction
	err = inUC.txUseCase.CommitTransaction(ctx, tx)
	if err != nil {
		inUC.log.Error("in usecase.CreateTaskHandle CommitTransaction error :%v", err)

		return taskID, err
	}
	//-----------------------------------

	return taskID, nil
}

func (inUC *InUseCase) UpdateTaskHandle(ctx context.Context, task *entity.Task) (int, error) {
	// resp, err := taskUC.dbRepo.UpdateDBTask(context.Background(), task)
	// if err != nil {
	// 	return 0, err
	// }
	return 0, nil
} // TODO

func (inUC *InUseCase) DeleteTaskHandle(ctx context.Context, taskID int) error {
	// err := taskUC.dbRepo.DeleteDBTask(context.Background(), taskId)
	// if err != nil {
	// 	return err
	// }
	return nil
} // TODO

func (inUC *InUseCase) GetTaskHandle(ctx context.Context, taskID int) (entity.Task, error) { // TODO
	// resp, err := taskUC.dbRepo.GetDBTask(context.Background(), taskId)
	// if err != nil {
	// 	return entity.Task{}, err
	// }
	return entity.Task{}, nil
} // TODO

func (inUC *InUseCase) ListTaskHandle(ctx context.Context) ([]entity.Task, error) {
	emptytaskList := make([]entity.Task, 0)
	// taskList, err := taskUC.dbRepo.ListDBTask(context.Background())
	// if err != nil {
	// 	return emptytaskList, err
	// }
	return emptytaskList, nil
} // TODO (with filtr)
