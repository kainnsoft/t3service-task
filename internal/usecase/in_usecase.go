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
	taskUseCase TaskUseCase
	userUseCase UserUseCase
	txUseCase   TxUseCase
	kafkaClient *repo.KafkaProducers
	log         *logging.ZeroLogger
}

func NewInUseCase(ti TaskDBRepoInterface,
	ta TaskApproversDBRepoInterface,
	te TaskEventsDBRepoInterface,
	ui UserDBRepoInterface,
	tx TxDBRepoInterface,
	kafkaClient *repo.KafkaProducers,
	log *logging.ZeroLogger) *InUseCase {

	taskUC := NewTaskUseCase(ti, ui, ta, te, log)
	userUC := NewUserUseCase(ui, ta, te)
	txUC := NewTxUseCase(tx)

	return &InUseCase{*taskUC, *userUC, *txUC, kafkaClient, log}
}

var _ v1.TaskHandlerInterface = (*InUseCase)(nil)

func (inUC *InUseCase) CreateTaskHandle(ctx context.Context, data []byte, authorEmail string) (int, error) {
	var taskID int

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

		var curApprover entity.User
		curApprover, err = inUC.userUseCase.CheckAndReturnUserByEmail(ctx, v.Email)
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
	err = inUC.userUseCase.InsertTaskApprovers(ctx, tx, taskID, receivedTask.Approvers)
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
	err = inUC.taskUseCase.InsertTaskEvent(ctx, taskID, receivedTask.Author.ID, entity.Created)
	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle inUC.taskUseCase.InsertTaskEvent RollbackTransaction error :%w", rberr)
		}
		err = fmt.Errorf("in usecase.CreateTaskHandle inUC.taskUseCase.InsertTaskEvent error: %w", err)

		return taskID, err
	}

	// ------------- kafka (start)
	// send kafka message to analytic service about task create
	// отправляем данные о свершившемся событии - task создан, согласован, обновлён и т.д.,
	// а также о юзере, который инициировал это событие (либо автор, либо согласующий)
	err = repo.SendMessagesToKafka(
		inUC.kafkaClient.KafProducerAboutTaskEvent,
		&receivedTask,
		entity.Created,            // тип события
		receivedTask.Author.Email, // user email
		entity.AboutTaskEvent)     // topic

	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle repo.SendMessagesToKafka (analityc) RollbackTransaction error :%v", rberr)
		}
		err = fmt.Errorf("in_usecase.CreateTaskHandle repo.SendMessagesToKafka (analityc) error: %v", err)

		return taskID, err
	}

	// send kafka message to mail service about task create
	// отправляем данные о том, что нужно сделать с task-ом - согласовать, зареджектить и т.д.,
	// а также о юзере, который должен это событие произвести
	err = repo.SendMessagesToKafka(
		inUC.kafkaClient.KafProducerToMailService,
		&receivedTask,
		entity.ToApprove,
		receivedTask.Approvers[0].Email,
		entity.ToMailService)

	if err != nil {
		rberr := inUC.txUseCase.RollbackTransaction(ctx, tx)
		if rberr != nil {
			inUC.log.Error("in usecase.CreateTaskHandle repo.SendMessagesToKafka (mail) RollbackTransaction error :%v", rberr)
		}
		err = fmt.Errorf("in_usecase.CreateTaskHandle repo.SendMessagesToKafka (mail) error: %v", err)

		return taskID, err
	}

	// ------------- kafka (end)

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
	return 0, nil
} // TODO

func (inUC *InUseCase) DeleteTaskHandle(ctx context.Context, taskID int) error {
	return nil
} // TODO

func (inUC *InUseCase) GetTaskHandle(ctx context.Context, taskID int) (entity.Task, error) {
	task, err := inUC.taskUseCase.GetOneTask(ctx, taskID)
	if err != nil {
		inUC.log.Error("in_usecase.GetTaskHandle inUC.taskUseCase.GetOneTask error :%w", err)

		return entity.Task{}, err
	}

	return task, nil
}

func (inUC *InUseCase) GetListTaskHandle(ctx context.Context) ([]entity.Task, error) {
	taskList, err := inUC.taskUseCase.GetTaskList(ctx)
	if err != nil {
		inUC.log.Error("in_usecase.ListTaskHandle inUC.taskUseCase.GetListTaskHandle error :%w", err)

		return nil, err
	}

	return taskList, nil
}
