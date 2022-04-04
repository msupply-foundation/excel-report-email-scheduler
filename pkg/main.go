package main

import (
	"os"

	"excel-report-email-scheduler/pkg/reportEmailer"

	"github.com/bugsnag/bugsnag-go"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/robfig/cron"
)

// Main entry point of the backend plugin.
// Using the provided SDK, it is required to simply make a call to
// datasource.Serve() with the options required (as below)
// {
//  QueryDataHandler: A datasource with a QueryData method which handles /tsdb/query calls.
//  CallResourceHandler: An HTTP handler. Create with sdk-go/backend/resource/httpadapter
//  CheckHealthHandler: A datasource with a CheckHealth method which handles /api/{pluginID}/health endpoint
//  TransformDataHandler: Transformer [experimental]
//  GRPCSettings: Settings..?
// }
func main() {

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "90618551d4e8d52a45260c61033093df",
		ReleaseStage: "production",
	})

	log.DefaultLogger.Info("Starting up")
	serveOptions, sql, err := getServeOptions()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Error when retreiving serveOptions", err.Error())
		panic(err)
	}

	bugsnag.OnBeforeNotify(func(e *bugsnag.Event, c *bugsnag.Configuration) error {
		settings, err := sql.GetSettings()
		if err != nil {
			e.MetaData.Add("Additional Meta Data", "Settings", "The database is not working - could not extract URL")
		} else {
			e.MetaData.Add("Additional Meta Data", "Settings", settings)
		}

		return nil
	})

	re := reportEmailer.NewReportEmailer(sql)

	// Try to send reports on loading
	re.CreateReports()

	// Set up scheduler which will try to send reports every 10 minutes
	c := cron.New()
	c.AddFunc("@every 10m", re.CreateReports)
	c.Start()

	// Start listening to requests sent from Grafana. This call is blocking and
	// and waits until Grafana shutsdown or the plugin exits.
	err = datasource.Serve(serveOptions)

	if err != nil {
		log.DefaultLogger.Error("FATAL. Plugin datasource.Serve() returned an error: " + err.Error())
		panic(err)
		os.Exit(1)
	}
}
