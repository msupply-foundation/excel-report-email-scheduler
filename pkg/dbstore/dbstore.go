package dbstore

import (
	"context"
	"database/sql"
	"runtime"

	"net/http"
	"os"
	"path/filepath"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	_ "modernc.org/sqlite"
)

// TODOs:
// SQL In separate files.
// Repository structs for each table? i.e. `ScheduleRepository` which has the methods `ScheduleRepository.update()` etc to partition the datasource struct.
// Common serialization/deserialization methods
// Better/consistent erroring
// Better/consistent return values

// Basic SQLite datasource
type SQLiteDatasource struct {
	instanceManager instancemgmt.InstanceManager
	Path            string
}

// Basic plugin instance settings
type InstanceSettings struct {
	httpClient *http.Client
}
type queryModel struct {
	QueryText   string   `json:"queryText"`
	TimeColumns []string `json:"timeColumns"`
}

func GetDataSource() *SQLiteDatasource {
	var runningOS string
	var dataPath string

	log.DefaultLogger.Info("GetDatasource")

	instanceManager := datasource.NewInstanceManager(getDataSourceInstanceSettings)

	runningOS = runtime.GOOS
	log.DefaultLogger.Info("OS Check: You are running %s platform.", runningOS)

	if runningOS == "windows" {
		dataPath = filepath.Join("..", "data", "msupply.db")
	} else {
		dataPath = filepath.Join("/var/lib/grafana/plugins", "data", "msupply.db")
	}

	sqlDatasource := &SQLiteDatasource{
		instanceManager: instanceManager,
		Path:            dataPath,
	}

	sqlDatasource.Init()

	return sqlDatasource
}

func getDataSourceInstanceSettings(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &InstanceSettings{
		httpClient: &http.Client{},
	}, nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (datasource *SQLiteDatasource) QueryData(ctx context.Context, request *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if os.Getenv("APP_ENV") != "production" {
		log.DefaultLogger.Info("QueryData", "request", request)
	}

	// TODO: Use context to timeout/cancel?
	response := backend.NewQueryDataResponse()

	for _, query := range request.Queries {
		res := datasource.query(ctx, query)
		response.Responses[query.RefID] = res
	}

	return response, nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (datasource *SQLiteDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	err := datasource.Ping()

	var status backend.HealthStatus
	var message string
	if err != nil {
		status = backend.HealthStatusError
		message = "Could not ping the database: " + err.Error()
	} else {
		status = backend.HealthStatusOk
		message = "Yeah, nah, All good"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (datasource *SQLiteDatasource) Ping() error {
	log.DefaultLogger.Info("Pinging Database")

	db, err := sql.Open("sqlite", datasource.Path)
	if err != nil {
		log.DefaultLogger.Error("Ping - sql.Open: ", err.Error())
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.DefaultLogger.Warn("Ping - sql.Ping: ", err.Error())
		return err
	}

	return nil
}

func (datasource *SQLiteDatasource) Init() {
	log.DefaultLogger.Info("Initializing Database")

	err := datasource.Ping()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Init - Ping. ", err.Error())
		panic(err)
	}

	db, err := sql.Open("sqlite", datasource.Path)
	if err != nil {
		log.DefaultLogger.Error("FATAL. Init - sql.Open : ", err.Error())
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS Schedule (id TEXT PRIMARY KEY, interval INTEGER, nextReportTime INTEGER, name TEXT, description TEXT, lookback INTEGER, reportGroupID TEXT, time TEXT, day INTEGER, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id))")
	stmt.Exec()
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Could not create Schedule:", err.Error())
		panic(err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS Config (id TEXT PRIMARY KEY, grafanaUsername TEXT, grafanaPassword TEXT, emailPassword TEXT, email TEXT, datasourceID INTEGER, emailHost TEXT, emailPort INTEGER, grafanaURL TEXT)")
	stmt.Exec()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Could not create Config:", err.Error())
		panic(err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportGroup (id TEXT PRIMARY KEY, name TEXT, description TEXT)")
	stmt.Exec()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Could not create ReportGroup:", err.Error())
		panic(err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportGroupMembership (id TEXT PRIMARY KEY, userID TEXT, reportGroupID TEXT, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id))")
	stmt.Exec()

	if err != nil {
		log.DefaultLogger.Error("FATAL. Could not create ReportGroupMembership:", err.Error())
		panic(err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportContent (id TEXT PRIMARY KEY, scheduleID TEXT, panelID INTEGER, dashboardID TEXT, lookback INTEGER, variables TEXT, FOREIGN KEY(scheduleID) REFERENCES Schedule(id))")
	stmt.Exec()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Could not create ReportContent:", err.Error())
		panic(err)
	}

	log.DefaultLogger.Info("Database initialized!")
}

// getDSInstance Returns cached datasource or creates new one
func (ds *SQLiteDatasource) GetDSInstance(pluginContext backend.PluginContext) (*InstanceSettings, error) {
	log.DefaultLogger.Debug("pluginContext.DataSourceInstanceSettings.ID", pluginContext.DataSourceInstanceSettings.ID)
	instance, err := ds.instanceManager.Get(pluginContext)
	if err != nil {
		log.DefaultLogger.Error("GetDSInstance", err.Error())
		return nil, err
	}
	return instance.(*InstanceSettings), nil
}
