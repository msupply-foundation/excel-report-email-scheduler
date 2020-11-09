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

func GetEmails(authConfig auth.AuthConfig, userIDs []string) []string {

	url := "http://" + "admin:admin@" + "localhost:3000/api/tsdb/query"

	queryString := "("
	i := 0
	for i < len(userIDs)-1 {
		queryString += "'" + userIDs[i] + "'" + ", "
		i += 1
	}
	queryString += "'" + userIDs[i] + "')"

	body, e := NewQueryRequest("SELECT * FROM \"user\" WHERE id IN "+queryString, "0", "0").ToRequestBody()

	if e != nil {
		log.DefaultLogger.Error(e.Error())
	}

	response, e := http.Post(url, "application/json", body)

	if e != nil {
		log.DefaultLogger.Error(e.Error())
	}

	qr, e := NewQueryResponse(response)

	if e != nil {
		log.DefaultLogger.Error(e.Error())
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

	return ids

}
