# mSupply Data Source

This data source is a backend datastore for the mSupply Dashboard App plugin.

The backend plugin uses a simple SQLite database and runs a HTTP restFUL server over Grafana's internal gRPC network to serve specific query requests.

Currently the mSupply Dashboard App Plugin supports automatic emailing of reports which this datasource holds configurations for.

Currently querying this datasource is unsupported for regular grafana operations.
