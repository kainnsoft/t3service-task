package usecase

import (
	"context"
	"testing"

	"team3-task/internal/entity"
	repo_mocks "team3-task/mocks/usecase"
	"team3-task/pkg/logging"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	Author    = entity.User{ID: 0, Email: "author@mail.ru"}
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

func TestGetOneTask(t *testing.T) {
	req := require.New(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	taskRepoMock := repo_mocks.NewMockTaskDBRepoInterface(mockCtrl)
	tu := repo_mocks.NewMockUserDBRepoInterface(mockCtrl)
	ta := repo_mocks.NewMockTaskApproversDBRepoInterface(mockCtrl)
	te := repo_mocks.NewMockTaskEventsDBRepoInterface(mockCtrl)
	log := logging.New("")

	l := NewTaskUseCase(taskRepoMock, tu, ta, te, log)

	t.Run("GetOneTask test", func(t *testing.T) {
		taskRepoMock.EXPECT().GetDBTask(context.Background(), taskID).Return(Task1, nil).Times(1)
		tu.EXPECT().GetDBUserByID(context.Background(), 0).Return(Author, nil).Times(1)
		ta.EXPECT().GetTaskApproversByTaskID(context.Background(), Task1.ID).Return(Task1.Approvers, nil).Times(1)

		task0, err := l.GetOneTask(context.Background(), taskID)

		req.NotNil(task0, "task must not be nil")
		req.Nil(err, "error must be nil")
	})
}
