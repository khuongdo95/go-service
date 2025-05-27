package tx

import (
	"context"

	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
)

func WithTransaction(ctx context.Context, client *ent.Client, fn Fn) *response.AppError {
	tx, err := client.Tx(ctx)
	if err != nil {
		return response.ServerError(err.Error())
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			if err := tx.Rollback(); err != nil {
				global.Log.Error("can not rollback transaction", err)
			}
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			if err := tx.Rollback(); err != nil {
				global.Log.Error("can not rollback transaction", err)
			}
		} else {
			// all good, commit
			if err := tx.Commit(); err != nil {
				global.Log.Error("can not rollback transaction", err)
			}
		}
	}()
	errApp := fn(ctx, tx)

	return errApp
}

// Tx is an interface that models the standard transaction in
// `ent/tx`.
//
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Tx interface {
	Client() *ent.Client
	OnRollback(f ent.RollbackHook)
	OnCommit(f ent.CommitHook)
}

type Fn func(context.Context, Tx) *response.AppError
