package auth

type AuthConfig struct {
	Username string
	Password string
	URL      string
}

type EmailConfig struct {
	Email    string
	Password string
	Host     string
	Port     int
}
