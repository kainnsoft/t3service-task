package repository

import (
	"context"
	"team3-task/internal/usecase"
	"team3-task/pkg/pg"

	"github.com/jackc/pgx/v4"
)

type TxDBRepo struct {
	*pg.DB
}

func NewTxDBRepo(pgdb *pg.DB) *TxDBRepo {
	return &TxDBRepo{pgdb}
}

var _ usecase.TxDBRepoInterface = (*TxDBRepo)(nil)

func (txRepo *TxDBRepo) BeginDBTransaction(ctx context.Context) (*pgx.Tx, error) {
	tx, err := txRepo.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (txRepo *TxDBRepo) CommitDBTransaction(ctx context.Context, txPtr *pgx.Tx) error {
	tx := *txPtr
	return tx.Commit(ctx)
}

func (txRepo *TxDBRepo) RollbackDBTransaction(ctx context.Context, txPtr *pgx.Tx) error {
	tx := *txPtr
	return tx.Rollback(ctx)
}
