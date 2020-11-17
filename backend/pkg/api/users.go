package api

import (
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
)

type User struct {
	ID    string `json:"ID"`
	Email string `json:"email"`
}

func GetEmails(authConfig auth.AuthConfig, userIDs []string, datasourceID int) ([]string, error) {
	url := "http://" + authConfig.AuthString() + authConfig.URL + "/api/tsdb/query"
	queryString := "("
	i := 0
	for i < len(userIDs)-1 {
		queryString += "'" + userIDs[i] + "'" + ", "
		i += 1
	}
	queryString += "'" + userIDs[i] + "')"
	body, err := NewQueryRequest("SELECT * FROM \"user\" WHERE id IN "+queryString, "0", "0", datasourceID).ToRequestBody()
	if err != nil {
		log.DefaultLogger.Error("GetEmails: NewQueryRequest: " + err.Error())
		return nil, err
	}

	response, err := http.Post(url, "application/json", body)
	if err != nil {
		log.DefaultLogger.Error("GetEmails: http.Post: " + err.Error())
		return nil, err
	}

	qr, err := NewQueryResponse(response)
	if err != nil {
		log.DefaultLogger.Error("GetEmails: NewQueryResponse: " + err.Error())
		return nil, err
	}

	emailColumnIdx := 0
	i = 0
	for _, column := range qr.Columns() {
		if column.Text == "e_mail" {
			emailColumnIdx = i
		}
		i += 1
	}

	var ids []string
	for _, row := range qr.Rows() {
		if str, ok := row[emailColumnIdx].(string); ok {
			ids = append(ids, str)
		}
	}

	return ids, nil

}
