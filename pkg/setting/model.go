package setting

type Settings struct {
	GrafanaUsername     string `json:"grafanaUsername"`
	GrafanaPassword     string `json:"grafanaPassword"`
	GrafanaURL          string `json:"grafanaURL"`
	SenderEmailAddress  string `json:"senderEmailAddress"`
	SenderEmailPassword string `json:"senderEmailPassword"`
	SenderEmailPort     int    `json:"senderEmailPort"`
	SenderEmailHost     string `json:"senderEmailHost"`
	DatasourceID        int    `json:"datasourceID"`
}
