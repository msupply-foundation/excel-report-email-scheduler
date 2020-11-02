package api

import (
	"net/http"

	"github.com/grafana/simple-datasource-backend/pkg/auth"
)

type Column struct {
	Text string `json:"text"`
}

type TablePanel struct {
	ID         int             `json:"id"`
	Title      string          `json:"title"`
	RawSql     string          `json:"rawSql"`
	Datasource string          `json:"datasource"`
	Rows       [][]interface{} `json:"rows"`
	Columns    []Column        `json:"columns"`
}

func NewTablePanel(id int, title string, rawSql string, datasource string) *TablePanel {
	return &TablePanel{ID: id, Title: title, RawSql: rawSql, Datasource: datasource}
}

func (panel *TablePanel) GetData(authConfig auth.AuthConfig) {
	body, _ := NewQueryRequest(panel.RawSql).ToRequestBody()

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
