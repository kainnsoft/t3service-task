package usecase

import (
	"context"
	"fmt"
	"team3-task/internal/entity"

	"github.com/jackc/pgx/v4"
)

type UserUseCase struct {
	userDBRepo UserDBRepoInterface
	taDBRepo   TaskApproversDBRepoInterface
	teDBRepo   TaskEventsDBRepoInterface
}

func NewUserUseCase(u UserDBRepoInterface,
	ta TaskApproversDBRepoInterface,
	te TaskEventsDBRepoInterface) *UserUseCase {

	return &UserUseCase{u, ta, te}
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

func (userUC *UserUseCase) InsertTaskApprovers(ctx context.Context, txp *pgx.Tx, taskId int, approverList []entity.User) error {
	err := userUC.taDBRepo.InsertDBTaskApprovers(ctx, txp, taskId, approverList)
	if err != nil {
		err = fmt.Errorf("usecase.InsertTaskApprovers userUC.taDBRepo.InsertDBTaskApprovers error: %v", err)

		return err
	}

	return nil
}

// получаем список всех approver-ов и всех task event-ов
// next - это первый из approver-ов по которому нет события в task event-е
func (userUC *UserUseCase) GetNextApprover(ctx context.Context, taskID int) (entity.User, error) {
	nextApprover := entity.User{}

	approverIDList, err := userUC.taDBRepo.GetTaskApproversIDByTaskID(ctx, taskID)
	if err != nil {
		err = fmt.Errorf("usecase.GetNextApprover userUC.taDBRepo.GetTaskApproversIDByTaskID error: %w", err)

		return nextApprover, err
	}

	eventApproversIDMap, err := userUC.teDBRepo.GetApproversIDMapMatchingTheListByTaskID(ctx, taskID, approverIDList)
	if err != nil {
		err = fmt.Errorf("usecase.GetNextApprover userUC.teDBRepo.GetApproversIDMatchingTheListByTaskID error: %w", err)

		return nextApprover, err
	}

	var nextApproverID int
	for _, v := range approverIDList {
		_, ok := eventApproversIDMap[v]
		if ok {
			continue
		}
		nextApproverID = v
		break
	}
	nextApprover, err = userUC.userDBRepo.GetDBUserByID(ctx, nextApproverID)
	if err != nil {
		err = fmt.Errorf("usecase.GetNextApprover userUC.userDBRepo.GetDBUserByID error: %w", err)

		return nextApprover, err
	}

	return nextApprover, nil
}
