package datasource

import (
	"context"
	"database/sql"
	"os"

	"excel-report-email-scheduler/pkg/ereserror"
)

type SqlClient struct {
	db  *sql.DB
	ctx *context.Context
	tx  *sql.Tx
}

func (datasource *MsupplyEresDatasource) NewSqlClient() (*SqlClient, error) {
	_, err := os.Stat(datasource.DataPath)
	if err != nil {
		err = ereserror.New(500, err, "datasource path does not exist")
		return nil, err
	}

	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		err = ereserror.New(500, err, "Failed to open datasource")
		return nil, err
	}

	return &SqlClient{db: db}, nil
}

func (client *SqlClient) BeginTx() error {
	ctx := context.Background()
	client.ctx = &ctx
	tx, err := client.db.BeginTx(ctx, nil)
	if err != nil {
		err = ereserror.New(500, err, "Begin transaction failed")
		return err
	}
	client.tx = tx
	return nil
}
