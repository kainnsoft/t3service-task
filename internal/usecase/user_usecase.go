package usecase

import (
	"context"
	"team3-task/internal/entity"
	"team3-task/pkg/logging"

	"github.com/jackc/pgx/v4"
)

type UserUseCase struct {
	userDBRepo UserDBRepoInterface
	log        *logging.ZeroLogger
}

func NewUserUseCase(r UserDBRepoInterface, log *logging.ZeroLogger) *UserUseCase {
	return &UserUseCase{r, log}
}

func (userUC *UserUseCase) CheckAndReturnUserByEmail(ctx context.Context, email string) (entity.User, error) {
	emptyUser := entity.User{}
	foundUser, err := userUC.GetUserByEmail(ctx, email)
	if (err != nil) && (err != pgx.ErrNoRows) {
		return emptyUser, err
	}

	if foundUser != emptyUser {
		return foundUser, nil
	}

	createdUser, err := userUC.CreateUser(ctx, email)
	if err != nil {
		return emptyUser, err
	}

	return createdUser, nil
}

func (userUC *UserUseCase) CreateUser(ctx context.Context, userEmail string) (entity.User, error) {
	emptyUser := entity.User{}
	_, err := userUC.userDBRepo.CreateDBUser(ctx, userEmail)
	if err != nil {
		return emptyUser, err
	}
	foundUser, err := userUC.userDBRepo.GetDBUserByEmail(ctx, userEmail)
	if err != nil {
		return emptyUser, err
	}
	return foundUser, nil
}

func (userUC *UserUseCase) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	emptyUser := entity.User{}
	resp, err := userUC.userDBRepo.GetDBUserByEmail(ctx, email)
	if err != nil {
		return emptyUser, err
	}
	return resp, nil
}
