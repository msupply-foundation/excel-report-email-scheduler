package api

import (
	"net/http"
	"strings"

	"github.com/grafana/simple-datasource-backend/pkg/auth"
)

type Column struct {
	Text string `json:"text"`
}

type TablePanel struct {
	ID           int             `json:"id"`
	From         string          `json:"from"`
	To           string          `json:"to"`
	Title        string          `json:"title"`
	RawSql       string          `json:"rawSql"`
	Rows         [][]interface{} `json:"rows"`
	Columns      []Column        `json:"columns"`
	Variables    TemplateList    `json:"variables"`
	DatasourceID int             `json:"DatasourceID"`
}

func NewTablePanel(id int, title string, rawSql string, from string, to string, datasourceID int) *TablePanel {
	return &TablePanel{ID: id, Title: title, RawSql: rawSql, From: from, To: to, DatasourceID: datasourceID}
}

func (panel *TablePanel) usesVariable(variable TemplateVariable) bool {
	return strings.Contains(panel.RawSql, "${"+variable.Name)
}

func (panel *TablePanel) injectVariable(variable TemplateVariable, storeIDs string) {

	if (variable.Name) == "store" {
		csv := ""
		split := strings.Split(storeIDs, ",")
		for i, substr := range split {

			if i == len(split)-1 {
				csv = csv + "'" + substr + "'"
			} else {
				csv = csv + "'" + substr + "'" + ", "
			}
		}

		panel.RawSql = strings.Replace(panel.RawSql, "${"+variable.Name+"}", csv, -1)

	} else {
		panel.RawSql = strings.Replace(panel.RawSql, "${"+variable.Name+"}", variable.Definition, -1)
	}
}

func (panel *TablePanel) PrepSql(variables TemplateList, storeIDs string) {
	for _, variable := range variables.List {
		if panel.usesVariable(variable) {
			panel.injectVariable(variable, storeIDs)
		}
	}
}

func (panel *TablePanel) GetData(authConfig auth.AuthConfig) {
	body, _ := NewQueryRequest(panel.RawSql, panel.From, panel.To).ToRequestBody()

	url := "http://" + authConfig.AuthString() + "localhost:3000/api/tsdb/query"
	response, _ := http.Post(url, "application/json", body)
	qr, _ := NewQueryResponse(response)

	panel.SetRows(qr.Rows())
	panel.SetColumns(qr.Columns())
}

func (panel *TablePanel) SetRows(rows [][]interface{}) {
	panel.Rows = rows
}

func (panel *TablePanel) SetColumns(columns []Column) {
	panel.Columns = columns
}
