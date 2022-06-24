package api

import (
	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/ereserror"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type MemberDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetMemberDeatailsFromUserIDs(authConfig *auth.AuthConfig, userIDs []string, datasourceID int) ([]MemberDetail, error) {
	frame := trace()
	authUrl, err := authConfig.AuthURL()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not retrive auth credentials")
		return nil, err
	}

	url := *authUrl + "/api/ds/query"

	queryString := "("
	i := 0
	for i < len(userIDs)-1 {
		queryString += "'" + userIDs[i] + "'" + ", "
		i += 1
	}
	queryString += "'" + userIDs[i] + "')"

	body, err := NewQueryRequest("SELECT id,name,e_mail FROM \"user\" WHERE id IN "+queryString, "0", "0", datasourceID).ToRequestBody()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not retrive user(s)")
		return nil, err
	}

	response, err := http.Post(url, "application/json", body)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not retrive user(s)")
		return nil, err
	}

	qr, err := NewQueryResponse(response)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not retrive user(s)")
		return nil, err
	}

	var members []MemberDetail
	for _, row := range qr.Rows() {
		member := MemberDetail{}

		member.ID = row[slices.IndexFunc(qr.Columns(), func(c Column) bool { return c.Text == "id" })].(string)
		member.Name = row[slices.IndexFunc(qr.Columns(), func(c Column) bool { return c.Text == "name" })].(string)
		member.Email = row[slices.IndexFunc(qr.Columns(), func(c Column) bool { return c.Text == "e_mail" })].(string)

		members = append(members, member)
	}

	return members, nil
}

func GetEmails(authConfig auth.AuthConfig, userIDs []string, datasourceID int) ([]string, error) {
	authURL, err := authConfig.AuthURL()
	if err != nil {
		log.DefaultLogger.Error("GetEmails: NewQueryRequest: " + err.Error())
		return nil, err
	}

	url := *authURL + "/api/ds/query"
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
