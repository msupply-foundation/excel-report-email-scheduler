package setting

import (
	"context"
	"excel-report-email-scheduler/pkg/ereserror"
	"runtime"

	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/pkg/errors"
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
	frame := trace()
	pluginCxt := httpadapter.PluginConfigFromContext(ctx)

	jsonData, err := simplejson.NewJson(pluginCxt.AppInstanceSettings.JSONData)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not retrive setting")
		return nil, err
	}

	grafanaUsername := jsonData.Get("grafanaUsername").MustString()
	grafanaURL := jsonData.Get("grafanaURL").MustString()
	senderEmailAddress := jsonData.Get("senderEmailAddress").MustString()
	senderEmailPort := jsonData.Get("senderEmailPort").MustInt()
	senderEmailHost := jsonData.Get("senderEmailHost").MustString()
	datasourceID := jsonData.Get("datasourceID").MustInt()

	var grafanaPassword string
	if securePassword, exists := pluginCxt.AppInstanceSettings.DecryptedSecureJSONData["grafanaPassword"]; exists {
		grafanaPassword = securePassword
	} else {
		// Fallback
		grafanaPassword = jsonData.Get("grafanaPassword").MustString()
	}

	var emailPassword string
	if secureEmailPassword, exists := pluginCxt.AppInstanceSettings.DecryptedSecureJSONData["senderEmailPassword"]; exists {
		emailPassword = secureEmailPassword
	} else {
		// Fallback
		emailPassword = jsonData.Get("senderEmailPassword").MustString()
	}

	return &Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, GrafanaURL: grafanaURL, SenderEmailAddress: senderEmailAddress, SenderEmailPort: senderEmailPort, SenderEmailPassword: emailPassword, SenderEmailHost: senderEmailHost, DatasourceID: datasourceID}, nil
}

func trace() *runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return &frame
}
