package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

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
				Columns []struct {
					Text string `json:"text"`
				} `json:"columns"`
				Rows [][]interface{} `json:"rows"`
			} `json:"tables"`
			Dataframes interface{} `json:"dataframes"`
		} `json:"A"`
	} `json:"results"`
}

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

// TODO: Need to get datasourceID from somewhere
func NewQueryRequest(rawSql string) *QueryRequest {
	query := Query{RawSQL: rawSql, DatasourceID: 1, Format: "table", RefID: "A"}
	queryRequest := QueryRequest{From: "0", To: "0", Queries: []Query{query}}

	return &queryRequest
}

func (qr *QueryRequest) ToRequestBody() (*strings.Reader, error) {
	parsed, err := json.Marshal(qr)

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	body := strings.NewReader(string(parsed))

	return body, nil
}

func QueryFromResponse(response *http.Response) (*QueryResponse, error) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	var qr QueryResponse
	err = json.Unmarshal(body, &qr)

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	return &qr, nil
}
