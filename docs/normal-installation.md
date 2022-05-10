# Normal development installation

The instructions below would walk you through installing the app plugin in your normal system without using docker. The plugin has docker-compose.yml, Dockerfile, config.env.example and config/grafana.ini files for docker installation, these files would have no effect in your installation.

## Pre-requisite

- Setup grafana for development [here](https://github.com/grafana/grafana/blob/master/contribute/developer-guide.md) - this includes having node, go and yarn.
  - You would need mSupply dashboard's postgres database installed in your system.
  - The plugin expects the mSupply dashboard postgres datasource to be setup and enabled.
  - Node version 16 is recommended
  - Yarn package manager is recommended
  - Golang version 1.18 is recommended
    - Mage is expected to be installed globally but you can download Mage executable, put it in the root folder and run the build
  - Grafana v8.4.4 image is being used by the docker container.
  - For development at least, we are expecting the plugins to run unsigned. (instruction below)
    - We want to run the plugin unsigned because we would be constantly editing code and debugging and signed plugin would not like that.
- Clone this repo
- Do the coding

#### Tips and TIL

- Grafana reads plugins from its data/plugins directory. You can point the plugin lookup to a different directory in the grafana custom.ini file (more [here](https://grafana.com/docs/grafana/latest/administration/configuration/)). You may need to create a `custom.ini` file.
- Let's say the plugin is cloned at `Users/omsupply/Documents/Github/msupply-dashboard-app`, then the plugin path configuration would be as following:

  ```
  ...
  [paths]
  plugins = "Users/omsupply/Documents/Github/msupply-dashboard-app"
  ...
  ```

- The plugin is not signed so you might need to add this also to your `custom.ini`:

  ```
  ...
  [plugins]
  allow_loading_unsigned_plugins = msupplyfoundation-excelreportemailscheduler-app
  ...
  ```

## Steps

- Enable mSupply dashboard's postgresql datasource
  - The plugin expect mSupply dashboard's postgresql datasource to be the source of it's users and panel data. Make sure it is installed and ready with appropriate data.
  - Make sure the datasource is started in mSupply dashboard and is available to select in the plugin's configuration page.
  - Select it as the datasource.
- Make sure you have plugin's SQLite database ready
  - The plugin uses a SQLite database named `msupply.db` to store its internal data, that is report_group and scheduler and its variable data.
  - `plugins/data` folder is where `msupply.db` would live. This is our plugin's noSQL database. A blank database at least, is expected to be in this path.
    > Warning: The path is hard coded so it cannot be changed.
    > Without `msupply.db` database the plugin would not work at all.
    > Make sure it is there at the path the plugin expects. The path is `../data` for windows and `/var/lib/grafana/plugins/data` for linux (which is this docker installation).
- Download and add necessary plugins in plugins folder (optional)
  - If you want to use Grafana Panels that uses mSupply-table, msupply-worldmap, msupply-regionmap, msupply-horizontal-bar and other plugins, please download them and put them in plugins folder.
  - These plugins (and everything in plugins folder) are picked up if you have configured the repo's folder as Grafana's plugins folder (instruction above). The directory hierarchy does not matter, grafana takes the folder as a plugin folder if it has plugin.json file.
- Run `mage clean && mage -v && yarn start` command to start the plugin in watch mode.
  - `mage clean` cleans (deletes) dist folder to have fresh start
  - `mage -v` Use `mage -v` to build for all platforms.
  - `yarn start` will build frontend javascript code (react) and puts it in dist folder, it also listens for changes.
  - After the above command the dist folder is ready.
    > Note: The above command set is comprised of 3 commands. First 2 are golang commands and last one runs react app in development mode. If you have built the golang app and only want to debug react app, you can skip the first two commands.
- Restart grafana
  - Whenever golang executable changes you would need to restart grafana to allow those changes to be picked up.
  - No need to restart if the changes are just Javascript related.
- Enable the plugin from plugin catalogue
  - When you open the Grafana at `http://localhost:3000` and login. You should see the plugin in the Plugin catalogue.
  - Select and enable it.
- Add complete plugin configurations and you are good to go.

The UI uses only `@grafana/ui` [components](https://grafana.com/docs/grafana/latest/packages_api/ui/). The docs are still a WIP, I think. Best to use the storybook.

Using [react-query](https://github.com/tannerlinsley/react-query) for async calls. It's nice.

## Other things

- `mage -v` builds for all platforms. This would create more than 5 executables, but you are only using one.
  - Do `mage build:windows` if you want to build for Windows
  - Do `mage build:linuxARM64` if you want to build for Linux
