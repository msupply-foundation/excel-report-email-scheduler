# mSupply Dashboard App

mSupply dashboard app is a grafana app plugin. Contained is a frontend and backend. The frontend is a plugin which can be enabled (once installed) through the standard plugins list. The backend is a datasource which can be created through the standard datasource list.

# Getting Started

TLDR: All you really need to do is have grafana all setup for development: With Go, Node and Yarn, clone the grafana repo, `yarn install --pure-lockfile` to build the frontend assets, `make run` compiles the backend then `yarn start` will allow you to connect to `http://localhost:3000` for grafana. The directory `data/plugins` is searched for plugins and you can install them through there. You still need to enable plugins and create data sources within the grafana UI, though. For the plugin, cd ./backend && yarn install && go run mage.go && yarn watch and for the frontned: cd ./frontend && yarn install && yarn watch

#### Prerequisites

- Setup grafana for development [here](https://github.com/grafana/grafana/blob/master/contribute/developer-guide.md) - this includes having node, go and yarn.
- Clone this repo
- Do the coding

#### Tips and TIL:

- Grafana reads plugins from its data/plugins directory. You can point the plugin lookup to a different directory in the grafana custom.ini file (more [here](https://grafana.com/docs/grafana/latest/administration/configuration/)). You may need to create a `custom.ini` file - but you can then add the following:

```
[paths]
plugins = "/Users/joshuagriffin/repos/grafana-plugins"
```

- The plugin is not signed so you might need to add this also to your `custom.ini`:

```
[plugins]
allow_loading_unsigned_plugins = msupply, msupply-datasource
```

- `yarn watch` in the frontend project will watch for changes - but at the moment I don't know how to make grafana hot-reload and listen for the changes. This means you can make UI changes with a simple refresh of the webpage. For backend changes you need to restart the grafana server.
- `yarn storybook` is amazing (in grafana repo)

#### Weird things, maybe?

- `go run mage.go` in the backend will build the project - I couldn't get it to build properly with `go build`
- Essentially this is two projects in one and they share `package.json` etc .. :shrug: I reckon there's a way to combine the two plugins to get built together but I'm not sure, yet.
- Need to restart the server if you make any changes to `plugin.json`'s

# Frontend

To build: cd ./frontend && yarn install && yarn start

The frontend is for setting up reporting groups and schedules. This frontend communicates with the backend through a RESTful server running within the backend/datasource plugin.

The UI uses only `@grafana/ui` [components](https://grafana.com/docs/grafana/latest/packages_api/ui/). The docs are still a WIP, I think. Best to use the storybook.

Using [react-query](https://github.com/tannerlinsley/react-query) for async calls. It's nice.

# Backend

To build: cd ./backend && yarn install && yarn start && go run mage.go

The backend primarily creates a static binary which grafana executes. This binary runs an http server which listens on `http://xxxx/api/plugins/msupply-datasource/resources` and the various end points are [here](https://github.com/openmsupply/msupply-dashboard-app/blob/869132fa53b41601bf9459a7c0ab00bdf8ec5476/backend/pkg/http_handler.go#L65-L91) _documentation to come ;-)_
