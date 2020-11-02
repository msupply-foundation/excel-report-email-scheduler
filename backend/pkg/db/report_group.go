package dbstore

import (
	"database/sql"

	"github.com/google/uuid"
)

type ReportGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewReportGroup(ID string, name string, description string) *ReportGroup {
	return &ReportGroup{ID: ID, Name: name, Description: description}
}

func (datasource *SQLiteDatasource) ReportGroupFromSchedule(schedule Schedule) *ReportGroup {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	row := db.QueryRow("SELECT * FROM ReportGroup WHERE scheduleID = ?", schedule.ID)

	var ID, name, description string
	row.Scan(&ID, &name, &description)
	return NewReportGroup(ID, name, description)
}

func (datasource *SQLiteDatasource) GetReportGroups() []ReportGroup {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var reportGroups []ReportGroup

	rows, _ := db.Query("SELECT * FROM ReportGroup")
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description string
		rows.Scan(&ID, &Name, &Description)
		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	return reportGroups
}

func (datasource *SQLiteDatasource) CreateReportGroup() (ReportGroup, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	reportGroup := ReportGroup{ID: uuid.New().String(), Name: "New report group", Description: ""}
	stmt, _ := db.Prepare("INSERT INTO ReportGroup (id, name, description) VALUES (?,?,?)")
	stmt.Exec(reportGroup.ID, reportGroup.Name, reportGroup.Description)
	defer stmt.Close()

	return reportGroup, nil
}

func (datasource *SQLiteDatasource) UpdateReportGroup(id string, reportGroup ReportGroup) ReportGroup {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE ReportGroup SET name = ?, description = ? where id = ?")
	stmt.Exec(reportGroup.Name, reportGroup.Description, id)
	defer stmt.Close()

	return reportGroup
}

func (datasource *SQLiteDatasource) DeleteReportGroup(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportGroup WHERE id = ?")
	stmt.Exec(id)
	stmt, _ = db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	stmt.Exec(id)

	defer stmt.Close()

	return true, nil
}
