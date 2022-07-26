package v1

import (
	"context"
	"team3-task/internal/entity"
)

// интерфейс работы с юзкейсами, который реализуется в файлах пакета usecase, а "дергается" в методах http controller-а (/internal/controller/http/v1/)
type TaskHandlerInterface interface {
	CreateTaskHandle(context.Context, []byte, string) (string, error)
	UpdateTaskHandle(context.Context, entity.Task) (int, error)
	DeleteTaskHandle(context.Context, int) error             // int - Task.id
	GetTaskHandle(context.Context, int) (entity.Task, error) // int - Task.id
	ListTaskHandle(context.Context) ([]entity.Task, error)   // need add filter
}
