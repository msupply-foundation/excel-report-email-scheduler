package api

import (
	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/ereserror"
	"net/http"

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
	url := authConfig.AuthURL() + "/api/ds/query"

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
