package core

import (
	entsql "entgo.io/ent/dialect/sql"

	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/database/mysql"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	_ "github.com/khuongdo95/go-service/internal/generated/ent/runtime"
)

type EntClient struct {
	client *ent.Client
}

func NewEntClient(sqlDb *mysql.Connection, cgf *mysql.SQLConfig) (*EntClient, *response.AppError) {
	var err *response.AppError
	if sqlDb == nil {
		err = response.ServerError("sqlDb is nil")
		return nil, err
	}
	drv := entsql.OpenDB(cgf.RDBMS, sqlDb.DB())
	return &EntClient{
		client: ent.NewClient(ent.Driver(drv)),
	}, nil
}

func (ent EntClient) Client() *ent.Client {
	return ent.client
}
