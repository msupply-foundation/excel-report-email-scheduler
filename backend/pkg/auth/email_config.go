package auth

import (
	"database/sql"

	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

type EmailConfig struct {
	Email    string
	Password string
	Host     string
	Port     int
}

func NewEmailConfig(datasource *dbstore.SQLiteDatasource) *EmailConfig {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var email, password, host string
	var port int

	row := db.QueryRow("SELECT email, emailPassword, emailHost, emailPort as password FROM Config")
	row.Scan(&email, &password, &host, &port)

	return &EmailConfig{Email: email, Password: password, Host: host, Port: port}
}
