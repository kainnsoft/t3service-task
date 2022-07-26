// методы работы с БД. Реализация интерфейсов, которые будут "дергаться" в usecase-ах (пока отключена БД)
package repository

import (
	"context"
	"encoding/json"
	"os"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/usecase"
	"team3-task/pkg/logging"
)

const mockDBPath string = "./mockDB.json"

type TaskMockRepo struct {
	mockDB *os.File
	logger *logging.ZeroLogger // TODO скорее всего не нужно - удалить
}

var _ usecase.TaskDBRepoInterface = (*TaskMockRepo)(nil)

func NewTaskMockRepo(logger *logging.ZeroLogger) (*TaskMockRepo, error) {
	mockDB, err := os.OpenFile("mockdb.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		return nil, errors.Newf("error opening mockDB file: %v", err)
	}
	return &TaskMockRepo{mockDB: mockDB}, nil
}

func (repo *TaskMockRepo) CreateDBTask(ctx context.Context, task entity.Task) (string, error) {

	sliceOfByteTask, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return "", errors.Newf("repository TaskMockRepo Create marshal error: %v", err)
	}

	err = os.WriteFile(mockDBPath, sliceOfByteTask, 0664)
	if err != nil {
		return "", err
	}

	return "1", nil
}

func (repo *TaskMockRepo) UpdateDBTask(ctx context.Context, task entity.Task) (int, error) {
	//repo.Pool.Exec()
	return 0, nil // TODO
}

func (repo *TaskMockRepo) DeleteDBTask(ctx context.Context, id int) error {
	//repo.Pool.Exec()
	return nil // TODO
}

func (repo *TaskMockRepo) GetDBTask(ctx context.Context, id int) (entity.Task, error) {
	//repo.Pool.Exec()
	return entity.Task{}, nil // TODO
}

func (repo *TaskMockRepo) ListDBTask(ctx context.Context) ([]entity.Task, error) {
	//repo.Pool.Exec()
	return []entity.Task{}, nil // TODO
}
