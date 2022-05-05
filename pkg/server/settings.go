package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/setting"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (server *HttpServer) updateSettings(rw http.ResponseWriter, request *http.Request) {
	frame := trace()
	var settings setting.Settings
	requestBody, err := request.GetBody()
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = json.Unmarshal(bodyAsBytes, &settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = server.db.CreateOrUpdateSettings(settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = json.NewEncoder(rw).Encode(settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	server.Success(rw, "Setting successfully deleted")
}
