package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
)

type DashboardResponse struct {
	Meta struct {
		Type                  string    `json:"type"`
		CanSave               bool      `json:"canSave"`
		CanEdit               bool      `json:"canEdit"`
		CanAdmin              bool      `json:"canAdmin"`
		CanStar               bool      `json:"canStar"`
		Slug                  string    `json:"slug"`
		URL                   string    `json:"url"`
		Expires               time.Time `json:"expires"`
		Created               time.Time `json:"created"`
		Updated               time.Time `json:"updated"`
		UpdatedBy             string    `json:"updatedBy"`
		CreatedBy             string    `json:"createdBy"`
		Version               int       `json:"version"`
		HasACL                bool      `json:"hasAcl"`
		IsFolder              bool      `json:"isFolder"`
		FolderID              int       `json:"folderId"`
		FolderTitle           string    `json:"folderTitle"`
		FolderURL             string    `json:"folderUrl"`
		Provisioned           bool      `json:"provisioned"`
		ProvisionedExternalID string    `json:"provisionedExternalId"`
	} `json:"meta"`
	Dashboard struct {
		Annotations struct {
			List []struct {
				BuiltIn    int    `json:"builtIn"`
				Datasource string `json:"datasource"`
				Enable     bool   `json:"enable"`
				Hide       bool   `json:"hide"`
				IconColor  string `json:"iconColor"`
				Name       string `json:"name"`
				Type       string `json:"type"`
			} `json:"list"`
		} `json:"annotations"`
		Editable     bool          `json:"editable"`
		GnetID       interface{}   `json:"gnetId"`
		GraphTooltip int           `json:"graphTooltip"`
		ID           int           `json:"id"`
		Links        []interface{} `json:"links"`
		Panels       []struct {
			Datasource  string `json:"datasource"`
			FieldConfig struct {
				Defaults struct {
					Custom struct {
						Align      interface{} `json:"align"`
						Filterable bool        `json:"filterable"`
					} `json:"custom"`
					Mappings   []interface{} `json:"mappings"`
					Thresholds struct {
						Mode  string `json:"mode"`
						Steps []struct {
							Color string      `json:"color"`
							Value interface{} `json:"value"`
						} `json:"steps"`
					} `json:"thresholds"`
				} `json:"defaults"`
				Overrides []interface{} `json:"overrides"`
			} `json:"fieldConfig"`
			GridPos struct {
				H int `json:"h"`
				W int `json:"w"`
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"gridPos"`
			ID      int `json:"id"`
			Options struct {
				ShowHeader bool `json:"showHeader"`
				SortBy     []struct {
					Desc        bool   `json:"desc"`
					DisplayName string `json:"displayName"`
				} `json:"sortBy"`
			} `json:"options"`
			PluginVersion string `json:"pluginVersion"`
			Targets       []struct {
				Format       string        `json:"format"`
				Group        []interface{} `json:"group"`
				MetricColumn string        `json:"metricColumn"`
				RawQuery     bool          `json:"rawQuery"`
				RawSQL       string        `json:"rawSql"`
				RefID        string        `json:"refId"`
				Select       [][]struct {
					Params []string `json:"params"`
					Type   string   `json:"type"`
				} `json:"select"`
				Table          string        `json:"table"`
				TimeColumn     string        `json:"timeColumn"`
				TimeColumnType string        `json:"timeColumnType"`
				Where          []interface{} `json:"where"`
			} `json:"targets"`
			TimeFrom        interface{}   `json:"timeFrom"`
			TimeShift       interface{}   `json:"timeShift"`
			Title           string        `json:"title"`
			Transformations []interface{} `json:"transformations"`
			Type            string        `json:"type"`
		} `json:"panels"`
		SchemaVersion int           `json:"schemaVersion"`
		Style         string        `json:"style"`
		Tags          []interface{} `json:"tags"`
		Templating    struct {
			List []interface{} `json:"list"`
		} `json:"templating"`
		Time struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Timepicker struct {
		} `json:"timepicker"`
		Timezone string `json:"timezone"`
		Title    string `json:"title"`
		UID      string `json:"uid"`
		Version  int    `json:"version"`
	} `json:"dashboard"`
}

type Dashboard struct {
	Panels []TablePanel `json:"panels"`
	UID    string       `json:"uid"`
}

func NewDashboardResponse(response *http.Response) (*DashboardResponse, error) {
	var dashboardResponse DashboardResponse
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return &dashboardResponse, err
	}

	err = json.Unmarshal(body, &dashboardResponse)

	return &dashboardResponse, err
}

func NewDashboard(authConfig *auth.AuthConfig, uuid string) (*Dashboard, error) {
	url := "http://" + authConfig.AuthString() + "localhost:3000/api/dashboards/uid/" + uuid
	response, err := http.Get(url)

	if err != nil {
		log.DefaultLogger.Error("NewDashboard: HTTP Request", err.Error())
		return nil, err
	}

	dashboardResponse, err := NewDashboardResponse(response)

	if err != nil {
		log.DefaultLogger.Error("NewDashboard: NewDashboardResponse", err.Error())
		return nil, err
	}

	var panels []TablePanel
	for _, panel := range dashboardResponse.Dashboard.Panels {
		newPanel := NewTablePanel(panel.ID, panel.Title, panel.Targets[0].RawSQL, panel.Datasource)
		panels = append(panels, *newPanel)
	}

	return &Dashboard{UID: dashboardResponse.Dashboard.UID, Panels: panels}, nil
}

func (dashboard *Dashboard) Panel(panelID int) (*TablePanel, error) {
	for _, panel := range dashboard.Panels {
		if panel.ID == panelID {
			return &panel, nil
		}
	}
	return nil, nil
}

func (resp *DashboardResponse) GetRawSQL(panelID int) string {
	for _, panel := range resp.Dashboard.Panels {
		if panel.ID == panelID {
			return panel.Targets[0].RawSQL
		}
	}
	return ""
}
