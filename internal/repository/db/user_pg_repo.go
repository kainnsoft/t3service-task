// методы работы с БД. Реализация интерфейсов, которые будут "дергаться" в usecase-ах
package repository

import (
	"context"
	"strings"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/usecase"
	"team3-task/pkg/pg"
)

type UserPGRepo struct {
	*pg.DB
}

var _ usecase.UserDBRepoInterface = (*UserPGRepo)(nil)

func NewUserPGRepo(pgdb *pg.DB) *UserPGRepo {
	return &UserPGRepo{pgdb}
}

func (repo *UserPGRepo) CreateDBUser(ctx context.Context, userEmail string) (string, error) {
	commandTag, err := repo.Pool.Exec(ctx,
		"INSERT INTO task.users (email) VALUES ($1);",
		strings.TrimSpace(userEmail))

	if err != nil {
		return commandTag.String(), errors.AddErrorContext(err, "user", " CreateDBUser insert error")
	}
	return commandTag.String(), nil // AffectedRows : 1
}

func (repo *UserPGRepo) UpdateDBUser(ctx context.Context, user entity.User) (int, error) {
	// repo.Pool.Exec()
	return 0, nil // TODO
}

func (repo *UserPGRepo) DeleteDBUser(ctx context.Context, id int) error {
	// repo.Pool.Exec()
	return nil // TODO
}

func (repo *UserPGRepo) GetDBUserByID(ctx context.Context, id int) (entity.User, error) {
	queryStr := "SELECT id, email FROM task.users WHERE id = $1;"
	row := repo.Pool.QueryRow(ctx, queryStr, id)

	foundUser := entity.User{}
	err := row.Scan(&foundUser.ID, &foundUser.Email)
	if err != nil {
		return foundUser, err
	}

	return foundUser, nil
}

func (repo *UserPGRepo) ListDBUser(ctx context.Context) ([]entity.User, error) {
	// repo.Pool.Exec()
	return []entity.User{}, nil // TODO
}

func (repo *UserPGRepo) GetDBUserByEmail(ctx context.Context, email string) (entity.User, error) {
	emptyUser := entity.User{}

	row := repo.Pool.QueryRow(ctx, "select id, email from task.users where email = $1;", email)

	foundUser := new(entity.User)
	err := row.Scan(&foundUser.ID, &foundUser.Email)
	if err != nil {
		return emptyUser, err
	}

	return *foundUser, nil
}
