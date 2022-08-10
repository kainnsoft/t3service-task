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

const mockUserDBPath string = "./mockUserDB.json"

type UserMockRepo struct {
	mockDB *os.File
}

var _ usecase.UserDBRepoInterface = (*UserMockRepo)(nil)

func NewUserMockRepo(logger *logging.ZeroLogger) (*UserMockRepo, error) {
	mockDB, err := os.OpenFile("mockdb.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o664)
	if err != nil {
		return nil, errors.Newf("error opening mockDB file: %v", err)
	}
	return &UserMockRepo{mockDB: mockDB}, nil
}

func (repo *UserMockRepo) CreateDBUser(ctx context.Context, userEmail string) (string, error) {
	sliceOfByteTask, err := json.MarshalIndent(userEmail, "", "  ")
	if err != nil {
		return "", errors.Newf("repository UserMockRepo Create marshal error: %v", err)
	}

	err = os.WriteFile(mockUserDBPath, sliceOfByteTask, 0o600)
	if err != nil {
		return "", err
	}

	return "1", nil
}

func (repo *UserMockRepo) UpdateDBUser(ctx context.Context, user entity.User) (int, error) {
	// repo.Pool.Exec()
	return 0, nil // TODO
}

func (repo *UserMockRepo) DeleteDBUser(ctx context.Context, id int) error {
	// repo.Pool.Exec()
	return nil // TODO
}

func (repo *UserMockRepo) GetDBUserByID(ctx context.Context, id int) (entity.User, error) {
	// repo.Pool.Exec()
	return entity.User{}, nil // TODO
}

func (repo *UserMockRepo) ListDBUser(ctx context.Context) ([]entity.User, error) {
	// repo.Pool.Exec()
	return []entity.User{}, nil // TODO
}

func (repo *UserMockRepo) GetDBUserByEmail(ctx context.Context, email string) (entity.User, error) {
	emptyUser := entity.User{}

	return emptyUser, nil
}
