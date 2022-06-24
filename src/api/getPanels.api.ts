import { SelectableValue } from '@grafana/data';
import { getBackendSrv } from '@grafana/runtime';
import { DashboardMeta, DashboardResponse, Panel, SelectableVariable, Variable } from 'types';
import { panelUsesUnsupportedMacro, panelUsesVariable } from 'utils/checkers.utils';
import { getDatasource } from './getDatasource.api';

export const getPanels = async (datasourceID: number): Promise<Panel[]> => {
  const dashboards = (await getDashboards()) ?? [];

  const datasource = await getDatasource(datasourceID);
  const { name: datasourceName = '' } = datasource ?? {};

  const panels: Panel[] = dashboards
    .filter(({ panels }) => panels?.length > 0)
    .map(({ panels, templating, uid }) => {
      const mappedPanels = panels
        .filter(({ type }) => type === 'table' || type === 'table-old' || type === 'msupplyfoundation-table')
        .map((rawPanel) => {
          const { targets } = rawPanel;
          const [target] = targets;
          const { rawSql } = target;

          const { list } = templating;

          // Want to filter out any panel which uses variables which aren't supported
          // Supported variables are custom and query types where the query type must
          // use the datasource specified in mSupply App Configuration.
          const unusableVariables = list.filter((variable: Variable) => {
            const { type, datasource: varDatasource } = variable;

            if (type === 'datasource' || type === 'adhoc') {
              return true;
            } else if (type === 'query') {
              // Note: Having a strange issue where templating.list.datasource never populates
              if (varDatasource === undefined) {
                return false;
              }

              if (typeof varDatasource === 'object' && varDatasource !== null) {
                return varDatasource.type.includes(datasourceName);
              }

              return varDatasource !== datasourceName;
            } else {
              return false;
            }
          });

          const usesUnusableVariables = unusableVariables.some((variable) => {
            const { name: variableName } = variable;
            return panelUsesVariable(rawSql, variableName);
          });

          let error = '';
          if (usesUnusableVariables) {
            error = 'This panel uses an unsupported variable.';
          } else if (panelUsesUnsupportedMacro(rawSql)) {
            error = 'This panel uses an unsupported macro.';
          }

          return { ...rawPanel, error };
        })
        .map(({ targets, description, title, id, type, error }) => {
          const [target] = targets;
          const { rawSql } = target;

          const dashboardID = uid;

          const { list } = templating;

          const variables =
            list?.reduce((acc: any, variable: Variable) => {
              if (panelUsesVariable(rawSql, variable.name)) {
                return [...acc, variable];
              } else {
                return acc;
              }
            }, []) ?? [];

          const mappedPanel = { error, rawSql, description, title, id, variables, dashboardID, type };

          return mappedPanel;
        });
      return mappedPanels;
    })
    .flat();

  return panels;
};

export const getDashboards = async () => {
  const dashboardMeta = await searchForDashboards();
  const dashboardResponses = await Promise.all<DashboardResponse>(dashboardMeta.map(({ uid }) => getDashboard(uid)));
  return dashboardResponses.map(({ dashboard }) => dashboard);
};

export const getDashboard = async (uuid: string): Promise<DashboardResponse> => {
  return getBackendSrv().get(`./api/dashboards/uid/${uuid}`);
};

export const searchForDashboards = async (): Promise<DashboardMeta[]> => {
  return getBackendSrv().get('./api/search');
};

export const refreshPanelOptions = async (
  variable: Variable,
  datasourceID: number
): Promise<Array<SelectableValue<SelectableVariable>>> => {
  const { definition, name } = variable;

  const optionsResponse = await getBackendSrv().post('/api/ds/query', {
    queries: [
      {
        datasourceId: datasourceID,
        rawSql: definition,
        format: 'table',
      },
    ],
  });

  const frames = optionsResponse.results.A.frames[0];

  const {
    data: { values },
  } = frames;

  const rows = values?.flat().map((datum: string) => {
    const selectableVariable = { name, value: datum } as SelectableVariable;
    const selectableValue = { label: datum, value: selectableVariable } as SelectableValue;

    return selectableValue;
  });

  return rows;
};
