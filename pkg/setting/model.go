package setting

type Settings struct {
	GrafanaUsername string `json:"grafanaUsername"`
	GrafanaPassword string `json:"grafanaPassword"`
	GrafanaURL      string `json:"grafanaURL"`
	Email           string `json:"senderEmailAddress"`
	EmailPassword   string `json:"senderEmailPassword"`
	EmailPort       int    `json:"senderEmailPort"`
	EmailHost       string `json:"senderEmailHost"`
	DatasourceID    int    `json:"datasourceID"`
}

func SettingsFieldatasource() string {
	return "\n{\n\tgrafanaUsername string" +
		"\n\tgrafanaPassword string" +
		"\n\tgrafanaURL string" +
		"\n\temail string\n}" +
		"\n\temailPassword string\n}" +
		"\n\temailPort int\n}" +
		"\n\temailHost string\n}" +
		"\n\tDatasourceID int\n}"
}
