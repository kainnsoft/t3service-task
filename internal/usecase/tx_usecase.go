package usecase

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type TxUseCase struct {
	txRepo TxDBRepoInterface
}

func NewTxUseCase(r TxDBRepoInterface) *TxUseCase {
	return &TxUseCase{r}
}

func (txr *TxUseCase) BeginTransaction(ctx context.Context) (*pgx.Tx, error) {
	tx, err := txr.txRepo.BeginDBTransaction(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (txr *TxUseCase) CommitTransaction(ctx context.Context, txPtr *pgx.Tx) error {
	return txr.txRepo.CommitDBTransaction(ctx, txPtr)
}

func (txr *TxUseCase) RollbackTransaction(ctx context.Context, txPtr *pgx.Tx) error {
	return txr.txRepo.RollbackDBTransaction(ctx, txPtr)
}
