package setting

import (
	"context"

	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

func SettingsFields() string {
	return "\n{\n\tgrafanaUsername string" +
		"\n\tgrafanaPassword string" +
		"\n\tgrafanaURL string" +
		"\n\tsenderEmailAddress string\n}" +
		"\n\tSenderEmailPassword string\n}" +
		"\n\tsenderEmailPort int\n}" +
		"\n\tsenderEmailHost string\n}" +
		"\n\tdatasourceID int\n}"
}

func NewSettings(ctx context.Context) (*Settings, error) {
	pluginCxt := httpadapter.PluginConfigFromContext(ctx)

	jsonData, err := simplejson.NewJson(pluginCxt.AppInstanceSettings.JSONData)
	if err != nil {
		log.DefaultLogger.Error("simplejson.NewJson:" + err.Error())
		panic(err)
	}

	grafanaUsername := jsonData.Get("grafanaUsername").MustString()
	grafanaURL := jsonData.Get("grafanaURL").MustString()
	senderEmailAddress := jsonData.Get("senderEmailAddress").MustString()
	senderEmailPort := jsonData.Get("senderEmailPort").MustInt()
	senderEmailHost := jsonData.Get("senderEmailHost").MustString()
	datasourceID := jsonData.Get("datasourceID").MustInt()

	log.DefaultLogger.Debug("Secure password: ", pluginCxt.AppInstanceSettings.DecryptedSecureJSONData["grafanaPassword"])
	var grafanaPassword string
	if securePassword, exists := pluginCxt.AppInstanceSettings.DecryptedSecureJSONData["grafanaPassword"]; exists {
		grafanaPassword = securePassword
	} else {
		// Fallback
		grafanaPassword = jsonData.Get("grafanaPassword").MustString()
	}

	log.DefaultLogger.Debug("Secure password: ", pluginCxt.AppInstanceSettings.DecryptedSecureJSONData["senderEmailPassword"])
	var emailPassword string
	if secureEmailPassword, exists := pluginCxt.AppInstanceSettings.DecryptedSecureJSONData["senderEmailPassword"]; exists {
		emailPassword = secureEmailPassword
	} else {
		// Fallback
		emailPassword = jsonData.Get("senderEmailPassword").MustString()
	}

	return &Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, GrafanaURL: grafanaURL, SenderEmailAddress: senderEmailAddress, SenderEmailPort: senderEmailPort, SenderEmailPassword: emailPassword, SenderEmailHost: senderEmailHost, DatasourceID: datasourceID}, nil
}
