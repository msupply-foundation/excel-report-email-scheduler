package auth

import (
	"database/sql"

	dbstore "github.com/grafana/simple-datasource-backend/pkg/db"
)

type EmailConfig struct {
	email    string
	password string
}

// TODO: Handle error cases and also might need to add additional
// fields i.e. SMTP etc
func NewEmailConfig(datasource *dbstore.SQLiteDatasource) EmailConfig {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var email, password string

	row := db.QueryRow("SELECT email, emailPassword as password FROM Config")
	row.Scan(&email, &password)

	return EmailConfig{email: email, password: password}
}
