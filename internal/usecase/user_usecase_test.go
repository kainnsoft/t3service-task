package usecase

import (
	"context"
	"team3-task/internal/entity"
	repo_mocks "team3-task/mocks/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	email      = "author@gmail.com"
	userAuthor = entity.User{ID: 1, Email: email}
)

func TestGetOneUser(t *testing.T) {
	req := require.New(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	u := repo_mocks.NewMockUserDBRepoInterface(mockCtrl)
	ua := repo_mocks.NewMockTaskApproversDBRepoInterface(mockCtrl)
	ue := repo_mocks.NewMockTaskEventsDBRepoInterface(mockCtrl)

	l := NewUserUseCase(u, ua, ue)

	t.Run("GetDBUserByEmail test", func(t *testing.T) {
		u.EXPECT().GetDBUserByEmail(context.Background(), email).Return(userAuthor, nil).Times(1)

		uAuthor, err := l.GetUserByEmail(context.Background(), email)

		req.NotNil(uAuthor, "uAuthor must not be nil")
		req.Nil(err, "error must be nil")
		req.Equal(email, uAuthor.Email, "wrong user email")
	})
}
