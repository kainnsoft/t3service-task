package usecase_test

import (
	"context"
	"testing"

	"team3-task/internal/entity"
	"team3-task/internal/usecase"

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

func (r *mockTaskDBRepo) CreateDBTask(ctx context.Context, task entity.Task) (string, error) {
	args := r.Called(task)
	arg0 := args.Get(0)
	if arg0 == "" {
		return "", args.Error(1)
	}
	return arg0.(string), args.Error(1)
}

func (r *mockTaskDBRepo) UpdateDBTask(ctx context.Context, task entity.Task) (int, error)
func (r *mockTaskDBRepo) DeleteDBTask(ctx context.Context, taskId int) error
func (r *mockTaskDBRepo) GetDBTask(ctx context.Context, taskId int) (entity.Task, error)
func (r *mockTaskDBRepo) ListDBTask(ctx context.Context) ([]entity.Task, error)

// tests
type unitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &unitTestSuite{})
}

func (s *unitTestSuite) TestCreateDBTask() {
	r := new(mockTaskDBRepo)

	r.On("CreateDBTask", Task1).Return("1", nil) // Added

	l := usecase.NewTaskUseCase(r)

	strAnswer, err := l.CreateTask(context.Background(), Task1)

	s.Nil(err, "error must be nil")
	//s.NotNil(strAnswer, "user must not be empty")
	s.Equal(strAnswer, "1", "created")

	r.AssertExpectations(s.T())
}
