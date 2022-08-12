package usecase_test

import (
	"context"
	"testing"

	"team3-task/internal/entity"
	"team3-task/internal/usecase"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	Author    = entity.User{Email: "author@mail.ru"}
	Approver1 = entity.User{Email: "appr1@mail.ru"}
	Approver2 = entity.User{Email: "appr2@mail.ru"}
	Approver3 = entity.User{Email: "appr3@mail.ru"}
	Approver4 = entity.User{Email: "appr4@mail.ru"}
	Task1     = entity.Task{
		Author:    Author,
		Descr:     "descr",
		Body:      "body",
		Approvers: []entity.User{Approver1, Approver2, Approver3, Approver4},
	}
)

// mocks
type mockTaskDBRepo struct {
	mock.Mock
}

func (r *mockTaskDBRepo) CreateDBTask(ctx context.Context, txPtr *pgx.Tx, task *entity.Task) (int, error) {
	args := r.Called(task)
	arg0 := args.Get(0)
	if arg0 == 0 {
		return 0, args.Error(1)
	}
	return arg0.(int), args.Error(1)
}

func (r *mockTaskDBRepo) UpdateDBTask(ctx context.Context, task *entity.Task) (int, error)
func (r *mockTaskDBRepo) DeleteDBTask(ctx context.Context, taskId int) error
func (r *mockTaskDBRepo) GetDBTask(ctx context.Context, taskId int) (entity.Task, error)
func (r *mockTaskDBRepo) GetListDBTask(ctx context.Context) ([]entity.Task, error)

type mockTaskApproversDBRepo struct {
	mock.Mock
}

func (r *mockTaskApproversDBRepo) InsertDBTaskApprovers(context.Context, *pgx.Tx, int, []entity.User) error
func (r *mockTaskApproversDBRepo) GetTaskApproversByTaskID(context.Context, int) ([]entity.User, error)
func (r *mockTaskApproversDBRepo) GetTaskApproversIDByTaskID(context.Context, int) ([]int, error)

type mockTaskEventsDBRepo struct {
	mock.Mock
}

func (r *mockTaskEventsDBRepo) InsertDBTaskEvents(ctx context.Context, taskID, userID int, taskEventType entity.KafkaTypes) error
func (r *mockTaskEventsDBRepo) GetTaskEventTypeByName(cctx context.Context, taskEventType entity.KafkaTypes) (int, error)
func (r *mockTaskEventsDBRepo) GetApproversIDMapMatchingTheListByTaskID(ctx context.Context, taskID int, listApproversID []int) (map[int]struct{}, error)

// tests
type unitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &unitTestSuite{})
}

func (s *unitTestSuite) TestCreateDBTask() {
	r := new(mockTaskDBRepo)
	ta := new(mockTaskApproversDBRepo)
	te := new(mockTaskEventsDBRepo)

	r.On("CreateDBTask", Task1).Return(1, nil) // Added

	l := usecase.NewTaskUseCase(r, nil, ta, te, nil)

	strAnswer, err := l.CreateTask(context.Background(), nil, &Task1)

	s.Nil(err, "error must be nil")
	//s.NotNil(strAnswer, "user must not be empty")
	s.Equal(strAnswer, 1, "created")

	r.AssertExpectations(s.T())
}
