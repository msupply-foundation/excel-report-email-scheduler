package dbstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ReportGroupFields() string {
	return "\n{\n\tid string" +
		"\n\tname string" +
		"\n\tdescription string\n}"
}

func NewReportGroup(ID string, name string, description string) *ReportGroup {
	return &ReportGroup{ID: ID, Name: name, Description: description}
}

func (datasource *SQLiteDatasource) ReportGroupFromSchedule(schedule Schedule) (*ReportGroup, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("ReportGroupFromSchedule: sql.Open", err.Error())
		return nil, err
	}

	row := db.QueryRow("SELECT * FROM ReportGroup WHERE ID = ?", schedule.ReportGroupID)

	var ID, name, description string
	err = row.Scan(&ID, &name, &description)
	if err != nil {
		log.DefaultLogger.Error("ReportGroupFromSchedule: rows.Scan(): ", err.Error())
		return nil, err
	}

	return NewReportGroup(ID, name, description), nil
}

func (datasource *SQLiteDatasource) GetReportGroups() ([]ReportGroup, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("GetReportGroups: sql.Open", err.Error())
		return nil, err
	}

	var reportGroups []ReportGroup

	rows, err := db.Query("SELECT * FROM ReportGroup")
	defer rows.Close()
	if err != nil {
		log.DefaultLogger.Error("GetReportGroups: db.Query(): ", err.Error())
		return nil, err
	}

	for rows.Next() {
		var ID, Name, Description string
		err = rows.Scan(&ID, &Name, &Description)
		if err != nil {
			log.DefaultLogger.Error("GetReportGroups: rows.Scan(): ", err.Error())
			return nil, err
		}

		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	return reportGroups, nil
}

func (datasource *SQLiteDatasource) CreateReportGroup() (*ReportGroup, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroup: sql.Open(): ", err.Error())
		return nil, err
	}

	reportGroup := ReportGroup{ID: uuid.New().String(), Name: "New report group", Description: ""}
	stmt, err := db.Prepare("INSERT INTO ReportGroup (id, name, description) VALUES (?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroup: db.Prepare(): ", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(reportGroup.ID, reportGroup.Name, reportGroup.Description)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroup: stmt.Exec(): ", err.Error())
		return nil, err
	}

	return &reportGroup, nil
}

func (datasource *SQLiteDatasource) UpdateReportGroup(id string, reportGroup ReportGroup) (*ReportGroup, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("UpdateReportGroup: sql.Open(): ", err.Error())
		return nil, err
	}

	stmt, err := db.Prepare("UPDATE ReportGroup SET name = ?, description = ? where id = ?")
	if err != nil {
		log.DefaultLogger.Error("UpdateReportGroup: db.Prepare(): ", err.Error())
		return nil, err
	}
	_, err = stmt.Exec(reportGroup.Name, reportGroup.Description, id)
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("UpdateReportGroup: stmt.Exec(): ", err.Error())
		return nil, err
	}

	return &reportGroup, nil
}

func (datasource *SQLiteDatasource) DeleteReportGroup(id string) error {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: sql.Open(): ", err.Error())
		return err
	}

	stmt, err := db.Prepare("DELETE FROM ReportGroup WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: db.Prepare(): ", err.Error())
		return err
	}
	stmt.Exec(id)

	stmt, err = db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	stmt.Exec(id)

	return nil
}
