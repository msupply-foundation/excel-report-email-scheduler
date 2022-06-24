package validation

import (
	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/ereserror"
	"runtime"

	"github.com/pkg/errors"
)

type Validation struct {
	datasource *datasource.MsupplyEresDatasource
	sqlClient  *datasource.SqlClient
}

func New(datasource *datasource.MsupplyEresDatasource) (*Validation, error) {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}

	return &Validation{datasource: datasource, sqlClient: sqlClient}, nil
}

func trace() *runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return &frame
}
