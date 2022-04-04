package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Query struct {
	RefID         string `json:"refId"`
	IntervalMs    int    `json:"intervalMs"`
	MaxDataPoints int    `json:"maxDataPoints"`
	DatasourceID  int    `json:"datasourceId"`
	RawSQL        string `json:"rawSql"`
	Format        string `json:"format"`
}

type QueryRequest struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Queries []Query `json:"queries"`
}

func NewQuery(rawSql string, datasource int) *Query {
	return &Query{RawSQL: rawSql, DatasourceID: datasource, Format: "table", RefID: "A"}
}

func NewQueryRequest(rawSql string, from string, to string, datasourceID int) *QueryRequest {
	query := NewQuery(rawSql, datasourceID)
	queryRequest := &QueryRequest{From: from, To: to, Queries: []Query{*query}}
	return queryRequest
}

func (qr *QueryRequest) ToRequestBody() (*strings.Reader, error) {
	parsed, err := json.Marshal(qr)
	if err != nil {
		log.DefaultLogger.Error("ToRequestBody: json.Marshal: " + err.Error())
		return nil, err
	}

	body := strings.NewReader(string(parsed))
	return body, nil
}

type QueryResponse struct {
	Results struct {
		A struct {
			RefID string `json:"refId"`
			Meta  struct {
				ExecutedQueryString string `json:"executedQueryString"`
				RowCount            int    `json:"rowCount"`
			} `json:"meta"`
			Series interface{} `json:"series"`
			Tables []struct {
				Columns []Column        `json:"columns"`
				Rows    [][]interface{} `json:"rows"`
			} `json:"tables"`
			Frames []struct {
				Schema struct {
					Fields []struct {
						Name string `json:"name"`
						Type string `json:"type"`
					} `json:"fields"`
				} `json:"schema"`
				Data struct {
					Values [][]interface{} `json:"values"`
				}
			} `json:"frames"`
		} `json:"A"`
	} `json:"results"`
}

func NewQueryResponse(response *http.Response) (*QueryResponse, error) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.DefaultLogger.Error("NewQueryResponse: ioutil.ReadAll: " + err.Error())
		return nil, err
	}
	log.DefaultLogger.Debug(fmt.Sprintf("NewQueryResponse: body fields: %+v\n", body))

	var qr QueryResponse
	err = json.Unmarshal(body, &qr)
	if err != nil {
		log.DefaultLogger.Error("NewQueryResponse: json.Unmarshal: " + err.Error())
		return nil, err
	}

	return &qr, nil
}

func (qr *QueryResponse) Rows() [][]interface{} {
	values := qr.Results.A.Frames[0].Data.Values
	if len(values) > 0 {
		columnCount := len(values)
		if columnCount > 0 {
			var rows = make([][]interface{}, len(values[0]))
			for rownum, _ := range rows {
				row := make([]interface{}, columnCount)
				for colnum, value := range values {
					row[colnum] = value[rownum]
				}
				rows[rownum] = row
			}

			return rows
		}
	}

	return nil
}

func (qr *QueryResponse) Columns() []Column {
	fields := qr.Results.A.Frames[0].Schema.Fields

	if len(fields) > 0 {
		columns := make([]Column, len(fields))
		for i, _ := range columns {
			var column Column
			column.Text = fields[i].Name
			columns[i] = column
		}
		return columns
	}

	return nil
}