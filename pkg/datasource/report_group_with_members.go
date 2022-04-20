package datasource

import "excel-report-email-scheduler/pkg/api"

type ReportGroupWithMembership struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Members     []api.MemberDetail `json:"members"`
}
