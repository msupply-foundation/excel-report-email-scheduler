package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/dbstore"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type User struct {
	ID    string `json:"ID"`
	Email string `json:"email"`
}

type UserDetails struct {
	ID    string `json:"ID"`
	name  string `json:"name"`
	Email string `json:"email"`
}

func GetMembersDetailFromGroup(authConfig *auth.AuthConfig, groupMembers []dbstore.ReportGroupMembership, datasourceID int) ([]UserDetails, error) {
	url := authConfig.AuthURL() + "/api/ds/query"
	queryString := "("
	i := 0
	for i < len(groupMembers)-1 {
		queryString += "'" + groupMembers[i].UserID + "'" + ", "
		i += 1
	}
	queryString += "'" + groupMembers[i].UserID + "')"

	body, err := NewQueryRequest("SELECT id,name,e_mail FROM \"user\" WHERE id IN "+queryString, "0", "0", datasourceID).ToRequestBody()
	if err != nil {
		log.DefaultLogger.Error("GetMembersDetailFromGroup: NewQueryRequest: " + err.Error())
		return nil, err
	}

	response, err := http.Post(url, "application/json", body)
	if err != nil {
		log.DefaultLogger.Error("GetMembersDetailFromGroup: http.Post: " + err.Error())
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.DefaultLogger.Error("GetMembersDetailFromGroup: ioutil.ReadAll: " + err.Error())
		return nil, err
	}

	var result []UserDetails
	if err := json.Unmarshal(responseData, &result); err != nil { // Parse []byte to the go struct pointer
		log.DefaultLogger.Error("GetMembersDetailFromGroup: json.Unmarshal: " + err.Error())
		return nil, err
	}

	return result, nil
}

func GetEmails(authConfig auth.AuthConfig, userIDs []string, datasourceID int) ([]string, error) {
	url := authConfig.AuthURL() + "/api/ds/query"
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
