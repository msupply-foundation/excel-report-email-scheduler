package api

import (
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// TODO: Use UID in URL, use config details for auth.
func GetDashboard(uid string) (*DashboardResponse, error) {
	response, err := http.Get("http://admin:admin@localhost:3000/api/dashboards/uid/ZpH7V3tMz")

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	dashboardResponse, err := DashboardFromResponse(response)

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	return dashboardResponse, nil
}
