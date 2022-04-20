package datasource

import (
	"database/sql"
	"excel-report-email-scheduler/pkg/api"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportGroupWithMembership struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Members     []api.MemberDetail `json:"members"`
}

func (datasource *MsupplyEresDatasource) DeleteReportGroupsWithMembers(id string) error {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: sql.Open(): ", err.Error())
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM ReportGroup WHERE id = ?")
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: db.Prepare(): ", err.Error())
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	stmt, err = db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: db.Prepare(): ", err.Error())
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return nil
}
