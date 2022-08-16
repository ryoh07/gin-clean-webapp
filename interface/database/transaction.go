package database

import (
	"context"

	"github.com/go-xorm/xorm"
	"github.com/ryoh07/gin-clean-webapp/transaction"
)

var txKey = struct{}{}

type tx struct {
	*xorm.Engine
}

func NewTransaction(db *xorm.Engine) transaction.Transaction {
	return &tx{db}
}

// トランザクション開始
func (t *tx) DoInTx(ctx context.Context, f func(ctx context.Context) error) error {

	// トランザクション開始
	session := t.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, &txKey, session)
	err = f(ctx)
	if err != nil {
		session.Rollback()
		return err
	}
	if err = session.Commit(); err != nil {
		session.Rollback()
		return err
	}
	return nil
}

func GetTx(ctx context.Context) *xorm.Session {
	session := ctx.Value(&txKey).(*xorm.Session)
	return session
}
