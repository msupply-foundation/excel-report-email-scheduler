package auth

import (
	"database/sql"

	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

type EmailConfig struct {
	Email    string
	Password string
}

// TODO: Handle error cases and also might need to add additional
// fields i.e. SMTP etc
func NewEmailConfig(datasource *dbstore.SQLiteDatasource) *EmailConfig {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var email, password string

	row := db.QueryRow("SELECT email, emailPassword as password FROM Config")
	row.Scan(&email, &password)

	return &EmailConfig{Email: email, Password: password}
}
