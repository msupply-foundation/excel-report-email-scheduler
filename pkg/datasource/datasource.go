package datasource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	_ "modernc.org/sqlite"
)

type MsupplyEresDatasource struct {
	im       instancemgmt.InstanceManager
	logger   log.Logger
	DataPath string
}

type MsupplyEresDatasourceInstance struct {
	logger log.Logger
}

func NewMsupplyEresDatasource() *MsupplyEresDatasource {
	var runningOS string
	var dataPath string

	logger := log.New()

	im := datasource.NewInstanceManager(newMsupplyEresDatasourceInstance)

	runningOS = runtime.GOOS
	logger.Info("OS Check: You are running %s platform.", runningOS)

	if runningOS == "windows" {
		dataPath = filepath.Join("..", "data", "msupply.db")
	} else {
		dataPath = filepath.Join("/var/lib/grafana/plugins", "data", "msupply.db")
	}

	mSupplyEresDatasource := &MsupplyEresDatasource{
		im:       im,
		logger:   logger,
		DataPath: dataPath,
	}

	_, err := mSupplyEresDatasource.Init()
	if err != nil {
		log.DefaultLogger.Error("Failed to initiate mSupplyEresDatasource", err)
	}

	return mSupplyEresDatasource
}

func newMsupplyEresDatasourceInstance(dsSettings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	logger := log.New()
	logger.Debug("Initializing new data source instance")

	return &MsupplyEresDatasourceInstance{
		logger: logger,
	}, nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (datasource *MsupplyEresDatasource) QueryData(ctx context.Context, request *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if os.Getenv("APP_ENV") != "production" {
		log.DefaultLogger.Info("QueryData", "request", request)
	}

	// TODO: Use context to timeout/cancel?
	response := backend.NewQueryDataResponse()

	// for _, query := range request.Queries {
	// 	res := datasource.query(ctx, query)
	// 	response.Responses[query.RefID] = res
	// }

	return response, nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (datasource *MsupplyEresDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

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

func (datasource *MsupplyEresDatasource) Ping() error {
	log.DefaultLogger.Info("Pinging Database")

	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = fmt.Errorf("Ping: NewSqlClient() : %w", err)
		return err
	}

	err = sqlClient.db.Ping()
	if err != nil {
		err = fmt.Errorf("Ping - sql.Ping : %w", err)
		return err
	}

	return nil
}

func (datasource *MsupplyEresDatasource) Init() (*bool, error) {
	datasource.logger.Info("Initializing Database")

	err := datasource.Ping()
	if err != nil {
		err = fmt.Errorf("FATAL. Init - Ping. %w", err)
		return nil, err
	}

	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = fmt.Errorf("FATAL. Init - sql.Open : %w", err)
		return nil, err
	}

	stmt, err := sqlClient.db.Prepare("CREATE TABLE IF NOT EXISTS Schedule (id TEXT PRIMARY KEY, interval INTEGER, nextReportTime INTEGER, name TEXT, description TEXT, lookback INTEGER, reportGroupID TEXT, time TEXT, day INTEGER, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id))")
	stmt.Exec()
	defer stmt.Close()
	if err != nil {
		err = fmt.Errorf("FATAL. Could not create Schedule: %w", err)
		return nil, err
	}

	stmt, err = sqlClient.db.Prepare("CREATE TABLE IF NOT EXISTS ReportGroup (id TEXT PRIMARY KEY, name TEXT, description TEXT)")
	stmt.Exec()
	if err != nil {
		err = fmt.Errorf("FATAL. Could not create ReportGroup: %w", err)
		return nil, err
	}

	stmt, err = sqlClient.db.Prepare("CREATE TABLE IF NOT EXISTS ReportGroupMembership (id TEXT PRIMARY KEY, userID TEXT, reportGroupID TEXT, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id))")
	stmt.Exec()

	if err != nil {
		err = fmt.Errorf("FATAL. Could not create ReportGroupMembership: %w", err)
		return nil, err
	}

	stmt, err = sqlClient.db.Prepare("CREATE TABLE IF NOT EXISTS ReportContent (id TEXT PRIMARY KEY, scheduleID TEXT, panelID INTEGER, dashboardID TEXT, lookback INTEGER, variables TEXT, FOREIGN KEY(scheduleID) REFERENCES Schedule(id))")
	stmt.Exec()
	if err != nil {
		err = fmt.Errorf("FATAL. Could not create ReportContent: %w", err)
		return nil, err
	}

	datasource.logger.Info("Database initialized!")

	status := true
	return &status, nil
}
