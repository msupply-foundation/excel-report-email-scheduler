package reportEmailer

import (
	"excel-report-email-scheduler/pkg/api"
	"excel-report-email-scheduler/pkg/datasource"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ReportEmailer struct {
	datasource *datasource.MsupplyEresDatasource
	inProgress bool
}

type Emailer struct {
	email    string
	password string
	host     string
	port     int
}

type Reporter struct {
	templatePath string
	reports      map[string]Report
}

type Report struct {
	id           string
	dateFormat   string
	name         string
	templatePath string
	file         *excelize.File
	sheets       []api.TablePanel
}
