# Docker development recommended

## Pre-requisite

- You would need mSupply dashboard's postgres database installed in your system.
- The plugin expects the mSupply dashboard postgres datasource to be setup and enabled.
- Docker must be installed in your system, also it expects docker-compose installed
- Node version 16 is recommended
- Yarn package manager for Node is recommended
- Golang version 1.18 is recommended
  - Mage is expected to be installed globally but you can download Mage executable, put it in the root folder and run the build
- Grafana v8.4.4 image is being used by the docker container.
- For development at least we are expecting the plugins to run unsigned. We have setup `GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS` in config.env for this purposed.
  - We want to run the plugin unsigned because we would be constantly editing code and debugging and signed plugin would not like that.

## Steps

- Make sure you have plugin's noSQL database ready
  - `plugins/data` folder is where msupply.db would live. This is our plugin's noSQL database. A blank database at least, is expected to be in this path.
    > Waring: Without `msupply.db` database the plugin would not work at all. Make sure it is there at the path the plugin expects. The path is `../data` for windows and `<plugins folder>/data` for linux (which is this docker installation).
- Download and add necessary plugins in plugins folder (optional)
  - If you want to use Grafana Panels that uses mSupply-table, msupply-worldmap, msupply-regionmap, msupply-horizontal-bar and other plugins, please download them and put them in plugins folder.
  - These plugins (and everything in plugins folder) are picked up by docker and added to Grafana plugin folder.
- Rename `config.env.example` to `config.env`
  - Add the settings you want Grafana to run with, including admin username and password.
- Run `mage clean && mage build:linuxARM64 && yarn start` command to start the plugin in watch mode.
  - `mage clean` cleans (deletes) dist folder to have fresh start
  - `mage build:linuxARM64` will build for linux which is the image we are using for docker. Alternatively use `mage -v` to build for all platforms.
  - `yarn start` will build frontend javascript code (react) and puts it in dist folder, it also listens for changes.
  - After the above command the dist folder is ready.
    > Note: The above command set is comprised of 3 commands. First 2 are golang commands and last one runs react app in development mode. If you have built the golang app and only want to debug react app, you can skip the first two commands.
- Start docker
  - In a new terminal window run `docker-compose up`
  - At first run this command will build the docker container.
  - It will sync the dist folder inside the container, adding it to grafana plugins folder.
  - It will set the config from `config/grafana.ini` and `config.env` and run grafana in development mode.
- Enable the plugin from plugin catalogue
  - When you open the Grafana at `http://localhost:3000` and login. You should see the plugin in the Plugin catalogue.
  - Select and enable it.
- Add complete plugin configurations and you are good to go.
