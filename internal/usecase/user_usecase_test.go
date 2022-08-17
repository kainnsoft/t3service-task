package usecase

import (
	"context"
	"team3-task/internal/entity"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	email      = "author@gmail.com"
	userAuthor = entity.User{ID: 1, Email: email}
)

type mockUserDBRepo struct {
	mock.Mock
}

func (u *mockUserDBRepo) CreateDBUser(ctx context.Context, email string) (string, error) {
	return "", nil
}
func (u *mockUserDBRepo) UpdateDBUser(ctx context.Context, user entity.User) (int, error) {
	return 0, nil
}
func (u *mockUserDBRepo) DeleteDBUser(ctx context.Context, userID int) error {
	return nil
}
func (u *mockUserDBRepo) GetDBUserByID(ctx context.Context, userID int) (entity.User, error) {
	return entity.User{}, nil
}
func (u *mockUserDBRepo) ListDBUser(ctx context.Context) ([]entity.User, error) {
	return nil, nil
}

func (u *mockUserDBRepo) GetDBUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := u.Called(email)
	arg0 := args.Get(0)
	if arg0 == nil {
		return entity.User{}, args.Error(1)
	}

	return arg0.(entity.User), args.Error(1)
}

type mockApproversDBRepo struct {
	mock.Mock
}

func (u *mockApproversDBRepo) InsertDBTaskApprovers(ctx context.Context, tx *pgx.Tx, taskID int, list []entity.User) error {
	return nil
}
func (u *mockApproversDBRepo) GetTaskApproversByTaskID(ctx context.Context, taskID int) ([]entity.User, error) {
	return nil, nil
}
func (u *mockApproversDBRepo) GetTaskApproversIDByTaskID(ctx context.Context, taskID int) ([]int, error) {
	return nil, nil
}

type mockEventsDBRepo struct {
	mock.Mock
}

func (u *mockEventsDBRepo) InsertDBTaskEvents(ctx context.Context, taskID, userID int, taskEventType entity.KafkaTypes) error {
	return nil
}
func (u *mockEventsDBRepo) GetTaskEventTypeByName(cctx context.Context, taskEventType entity.KafkaTypes) (int, error) {
	return 0, nil
}
func (u *mockEventsDBRepo) GetApproversIDMapMatchingTheListByTaskID(ctx context.Context, taskID int, listApproversID []int) (map[int]struct{}, error) {
	return nil, nil
}

func TestGetOneUser(t *testing.T) {
	u := new(mockUserDBRepo)
	ua := new(mockApproversDBRepo)
	ue := new(mockEventsDBRepo)

	u.On("GetDBUserByEmail", email).Return(userAuthor, nil)

	l := NewUserUseCase(u, ua, ue)

	uAuthor, err := l.GetUserByEmail(context.Background(), email)

	req := require.New(t)
	req.NotNil(uAuthor, "uAuthor must not be nil")
	req.Nil(err, "error must be nil")
	req.Equal(email, uAuthor.Email, "wrong user email")

	u.AssertExpectations(t)
}
