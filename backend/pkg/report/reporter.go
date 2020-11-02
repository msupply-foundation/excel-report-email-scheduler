package report

import (
	dbstore "github.com/grafana/simple-datasource-backend/pkg/db"
)

type Reporter struct {
	datasource *dbstore.SQLiteDatasource
}

func NewReporter(sqlite *dbstore.SQLiteDatasource) *Reporter {
	return &Reporter{datasource: sqlite}
}
