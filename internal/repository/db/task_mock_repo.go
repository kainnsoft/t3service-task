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

	"github.com/jackc/pgx/v4"
)

const mockDBPath string = "./mockDB.json"

type TaskMockRepo struct {
	mockDB *os.File
}

var _ usecase.TaskDBRepoInterface = (*TaskMockRepo)(nil)

func NewTaskMockRepo(logger *logging.ZeroLogger) (*TaskMockRepo, error) {
	mockDB, err := os.OpenFile("mockdb.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o664)
	if err != nil {
		return nil, errors.Newf("error opening mockDB file: %v", err)
	}
	return &TaskMockRepo{mockDB: mockDB}, nil
}

func (repo *TaskMockRepo) CreateDBTask(ctx context.Context, txPtr *pgx.Tx, task *entity.Task) (int, error) {

	sliceOfByteTask, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return 0, errors.Newf("repository TaskMockRepo Create marshal error: %v", err)
	}

	err = os.WriteFile(mockDBPath, sliceOfByteTask, 0o600)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (repo *TaskMockRepo) UpdateDBTask(ctx context.Context, task *entity.Task) (int, error) {
	// repo.Pool.Exec()
	return 0, nil // TODO
}

func (repo *TaskMockRepo) DeleteDBTask(ctx context.Context, id int) error {
	// repo.Pool.Exec()
	return nil // TODO
}

func (repo *TaskMockRepo) GetDBTask(ctx context.Context, id int) (entity.Task, error) {
	// repo.Pool.Exec()
	return entity.Task{}, nil // TODO
}

func (repo *TaskMockRepo) GetListDBTask(ctx context.Context) ([]entity.Task, error) {
	// repo.Pool.Exec()
	return []entity.Task{}, nil // TODO
}
