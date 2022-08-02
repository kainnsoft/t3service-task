package usecase

import (
	"context"
	"team3-task/internal/entity"

	"github.com/jackc/pgx/v4"
)

// интерфейс работы с БД, которыq реализуtтся в файлах репозитория, а "дергаtтся" в файлах этого пакета (usecase)
type TaskDBRepoInterface interface {
	CreateDBTask(context.Context, entity.Task) (int, error)
	UpdateDBTask(context.Context, entity.Task) (int, error)
	DeleteDBTask(context.Context, int) error             // int - Task.id
	GetDBTask(context.Context, int) (entity.Task, error) // int - Task.id
	ListDBTask(context.Context) ([]entity.Task, error)   // need add filter
	//------------------------------------------------
	BeginTransaction(context.Context) (*pgx.Tx, error)
	CommitTransaction(context.Context, *pgx.Tx) error
	RollbackTransaction(context.Context, *pgx.Tx)
	//------------------------------------------------
}

type TaskApproversDBRepoInterface interface {
	InsertDBTaskApprovers(context.Context, *pgx.Tx, int, []entity.User) error // batch
	GetTaskApproversByTaskID(context.Context, int) ([]entity.User, error)
}

type UserDBRepoInterface interface {
	CreateDBUser(context.Context, string) (string, error)
	UpdateDBUser(context.Context, entity.User) (int, error)
	DeleteDBUser(context.Context, int) error
	GetDBUser(context.Context, int) (entity.User, error)
	GetDBUserByEmail(context.Context, string) (entity.User, error)
	ListDBUser(context.Context) ([]entity.User, error) // need add filter
}
