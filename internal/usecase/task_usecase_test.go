package usecase

import (
	"context"
	"testing"

	"team3-task/internal/entity"
	"team3-task/pkg/logging"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	Task0 = entity.Task{
		Finished: false,
		Author:   Author,
		Descr:    "descr",
		Body:     "body",
	}

	taskID = 1
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

func (r *mockTaskDBRepo) UpdateDBTask(ctx context.Context, task *entity.Task) (int, error) {
	return 0, nil
}
func (r *mockTaskDBRepo) DeleteDBTask(ctx context.Context, taskId int) error { return nil }
func (r *mockTaskDBRepo) GetDBTask(ctx context.Context, taskId int) (entity.Task, error) {
	return entity.Task{}, nil
}
func (r *mockTaskDBRepo) GetListDBTask(ctx context.Context) ([]entity.Task, error) {
	return nil, nil
}

// tests
func TestGetOneTask(t *testing.T) {
	req := require.New(t)

	r := new(mockTaskDBRepo)
	tu := new(mockUserDBRepo)
	ta := new(mockApproversDBRepo)
	te := new(mockEventsDBRepo)
	log := logging.New("")

	r.On("GetDBTask", taskID).Return(Task1, nil)

	l := NewTaskUseCase(r, tu, ta, te, log)

	t.Run("GetOneTask test", func(t *testing.T) {
		task0, _ := l.GetOneTask(context.Background(), taskID)

		req.NotNil(task0, "task must not be nil")
	})
}
