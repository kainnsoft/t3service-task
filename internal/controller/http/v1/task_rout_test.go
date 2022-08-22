package v1

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	_ "team3-task/internal/docs"
	"team3-task/internal/entity"
	auth_mocks "team3-task/mocks/app"
	ctrl_mocks "team3-task/mocks/controller/http/v1"
	"team3-task/pkg/logging"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_taskRoutes_checkValidation(t *testing.T) {
	type args struct {
		r *http.Request
	}

	type tRout struct {
		tHandler *ctrl_mocks.MockTaskHandlerInterface
		grpcCl   *auth_mocks.MockAuthAccessChecker
	}

	tests := []struct {
		name    string
		prepare func(tr *tRout)
		args    args
		want    entity.AuthResponse
		wantErr bool
	}{
		{
			name: "",
			prepare: func(tr *tRout) {
				tr.tHandler.EXPECT().CreateTaskHandle(context.Background(), gomock.Any(), gomock.Any()).Return(gomock.Any(), nil).Times(1)
				tr.tHandler.EXPECT().UpdateTaskHandle(context.Background(), gomock.Any()).Return(gomock.Any(), nil).Times(1)
				tr.tHandler.EXPECT().DeleteTaskHandle(context.Background(), gomock.Any()).Return(gomock.Any()).Times(1)
				tr.tHandler.EXPECT().GetTaskHandle(context.Background(), gomock.Any()).Return(gomock.Any(), nil).Times(1)
				tr.tHandler.EXPECT().GetListTaskHandle(context.Background()).Return(gomock.Any(), nil).Times(1)

				tr.grpcCl.EXPECT().CheckAccess(gomock.Any()).Return(entity.AuthResponse{}, fmt.Errorf("validation error occures while cookie getting")).Times(1)
			},
			args:    args{r: new(http.Request)},
			want:    entity.AuthResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			log := logging.New("")
			th := ctrl_mocks.NewMockTaskHandlerInterface(mockCtrl) // taskHandler
			gc := auth_mocks.NewMockAuthAccessChecker(mockCtrl)    // grpcClient

			tr := &taskRoutes{log, th, gc}

			// test
			got, err := tr.checkValidation(tt.args.r)

			if (err != nil) != tt.wantErr {
				t.Errorf("taskRoutes.checkValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskRoutes.checkValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}
