# mSupply Dashboard: Excel report e-mail scheduler

Grafana App plugin for mSupply Dashboard application built with Golang as backend and react as frontend.

The plugin takes data from panels of mSupply dashboard to generate excel reports. The reports are then emailed to a custom user group created with mSupply users pulled from mSupply Dashboard's datasource. The timing of the scheduler can be set in the plugin.

## Docker development recommended build

### Pre-requisite

- You would need mSupply dashboard's postgres database installed in your system.
- The plugin expects the mSupply dashboard postgres datasource to be setup and enabled.
- Docker must be installed in your system, also it expects docker-compose installed
- Node version 16 is recommended
- Yarn package manager for Node is recommended
- Golang version 1.18 is recommended
  - Mage is expected to be installed globally but you can download Mage executable, put it in the root folder and run the build
- Grafana v8.4.4 image is being used by the docker container

### Steps

- Rename `config.env.example` to `config.env` and add the settings you want Grafana to run with, including admin username and password.
- `mage clean && mage build:linuxARM64 && yarn start`
  - `mage clean` cleans dist folder to have fresh start
  - `mage build:linuxARM64` will build for linux which is the image we are using for docker. Alternatively use `mage -v` to build for all platforms.
  - `yarn start` will build frontend javascript code (react) and puts it in dist folder, it also listens for changes.
  - After the above command the dist folder is ready.
- Run `docker-compose up`
  - At first run this command will build the docker container.
  - It will sync the dist folder inside the container.
  - It will set the config from `config/grafana.ini` and `config.env` and run grafana in development mode.
- When you open the Grafana at `http://localhost:3000` and login. You should see the plugin in the Plugin catalogue.
  - Select and enable it.
- Add complete plugin configuration, instructions are on the UI.

## Production build

- Instructions coming soon
