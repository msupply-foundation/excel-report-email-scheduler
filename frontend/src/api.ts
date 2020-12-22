import { SelectableValue } from '@grafana/data';
import { panelUsesVariable, panelUsesUnsupportedMacro } from './common/utils/checkers';
import { DashboardResponse, DashboardMeta, CreateContentVars, ReportContent, SelectableVariable } from './common/types';
import { getBackendSrv } from '@grafana/runtime';
import { Variable, Panel, ReportGroupMember, Schedule, Store, ReportGroup, User } from 'common/types';

export const refreshPanelOptions = async (
  variable: Variable,
  datasourceID: number
): Promise<Array<SelectableValue<SelectableVariable>>> => {
  const { definition, name } = variable;
  const optionsResponse = await getBackendSrv().post('/api/tsdb/query', {
    queries: [
      {
        datasourceId: datasourceID,
        rawSql: definition,
        format: 'table',
      },
    ],
  });

  const rows = optionsResponse.results?.A?.tables?.[0]?.rows?.flat().map((datum: string) => {
    const selectableVariable = { name, value: datum } as SelectableVariable;
    const selectableValue = { label: datum, value: selectableVariable } as SelectableValue;

    return selectableValue;
  });

  return rows;
};

export const sendTestEmail = (scheduleID: string) =>
  getBackendSrv().get(`api/plugins/msupply-datasource/resources/test-email?schedule-id=${scheduleID}`);

export const getGroupMembers = (_: string, groupId: string): Promise<ReportGroupMember[]> =>
  getBackendSrv().get(`api/plugins/msupply-datasource/resources/report-group-membership/?group-id=${groupId}`);

export const getReportGroups = (): Promise<ReportGroup[]> =>
  getBackendSrv().get('api/plugins/msupply-datasource/resources/report-group');

export const getUsers = (datasourceID: number): Promise<User[]> => {
  return getBackendSrv()
    .post('/api/tsdb/query', {
      queries: [
        {
          datasourceId: datasourceID,
          rawSql: 'SELECT id, name, first_name, last_name, e_mail FROM "user"',
          format: 'table',
        },
      ],
    })
    .then(result => {
      const {
        results: {
          A: {
            tables: [{ rows, columns }],
          },
        },
      } = result;

      const columnsToExtract = ['id', 'name', 'first_name', 'last_name', 'e_mail'];
      const indexes = columns.reduce((acc: number[], { text }: any, i: number) => {
        if (columnsToExtract.includes(text)) {
          return [...acc, i];
        }
        return acc;
      }, []);

      return rows.map((rowData: any) => {
        return indexes.reduce((acc: any, value: any, i: any) => {
          return { ...acc, [columns[value].text]: rowData[i] };
        }, {});
      });
    });
};

export const getStores = (datasourceID: number): Promise<Store[]> => {
  return getBackendSrv()
    .post('/api/tsdb/query', {
      queries: [
        {
          datasourceId: datasourceID,
          rawSql: 'SELECT id, name, code FROM "store"',
          format: 'table',
        },
      ],
    })
    .then(result => {
      const {
        results: {
          A: {
            tables: [{ rows, columns }],
          },
        },
      } = result;

      const columnsToExtract = ['id', 'name', 'code'];
      const indexes = columns.reduce((acc: number[], { text }: any, i: number) => {
        if (columnsToExtract.includes(text)) {
          return [...acc, i];
        }
        return acc;
      }, []);

      return rows.map((rowData: any) => {
        return indexes.reduce((acc: any, value: any, i: any) => {
          return { ...acc, [columns[value].text]: rowData[i] };
        }, {});
      });
    });
};

export const createReportGroupMembership = async (params: any) => {
  const { user, reportGroupID } = params;

  const membership = { userID: user.id, reportGroupID };
  return getBackendSrv().post('./api/plugins/msupply-datasource/resources/report-group-membership', [membership]);
};

export const deleteReportGroupMembership = async (reportGroupMembership: ReportGroupMember) => {
  const { id } = reportGroupMembership;
  return getBackendSrv().delete(`./api/plugins/msupply-datasource/resources/report-group-membership/${id}`);
};

export const updateReportGroup = async (reportGroup: any) => {
  return getBackendSrv().put(`./api/plugins/msupply-datasource/resources/report-group/${reportGroup?.id}`, reportGroup);
};

export const deleteReportGroup = async (reportGroup: ReportGroup) => {
  return getBackendSrv().delete(`./api/plugins/msupply-datasource/resources/report-group/${reportGroup?.id}`);
};

export const createReportGroup = async () => {
  return getBackendSrv().post('./api/plugins/msupply-datasource/resources/report-group');
};

export const getSchedules = async () => {
  return getBackendSrv().get('./api/plugins/msupply-datasource/resources/schedule');
};

export const createSchedule = async () => {
  return getBackendSrv().post('./api/plugins/msupply-datasource/resources/schedule');
};

export const deleteSchedule = async (schedule: Schedule) => {
  return getBackendSrv().delete(`./api/plugins/msupply-datasource/resources/schedule/${schedule?.id}`);
};

export const getReportContent = async (_: string, scheduleID: string) => {
  const content: ReportContent[] = await getBackendSrv().get(
    `./api/plugins/msupply-datasource/resources/report-content?schedule-id=${scheduleID}`
  );
  return content;
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

export const getDatasource = async (id: number) => {
  return getBackendSrv().get(`./api/datasources/${id}`);
};

export const getPanels = async (datasourceID: number): Promise<Panel[]> => {
  const dashboards = (await getDashboards()) ?? [];
  const datasource = await getDatasource(datasourceID);
  const { name: datasourceName = '' } = datasource ?? {};

  const panels: Panel[] = dashboards
    .filter(({ panels }) => panels?.length > 0)
    .map(({ panels, templating, uid }) => {
      const mappedPanels = panels
        .filter(({ type }) => type === 'table' || type === 'table-old')
        .map(rawPanel => {
          const { targets } = rawPanel;
          const [target] = targets;
          const { rawSql } = target;

          const { list } = templating;

          // Want to filter out any panel which uses variables which aren't supported
          // Supported variables are custom and query types where the query type must
          // use the datasource specified in mSupply App Configuration.
          const unusableVariables = list.filter((variable: Variable) => {
            const { datasource, type } = variable;
            if (type === 'datasource' || type === 'adhoc') {
              return true;
            } else if (type === 'query') {
              return datasourceName !== datasource;
            } else {
              return false;
            }
          });

          const usesUnusableVariables = unusableVariables.some(variable => {
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

export const updateSchedule = async (schedule: Schedule) => {
  return getBackendSrv().put(`./api/plugins/msupply-datasource/resources/schedule/${schedule?.id}`, schedule);
};

export const createReportContent = async (contents: CreateContentVars) => {
  return getBackendSrv().post(`./api/plugins/msupply-datasource/resources/report-content`, contents);
};

export const deleteReportContent = async (contents: any) => {
  return getBackendSrv().delete(`./api/plugins/msupply-datasource/resources/report-content/${contents?.id}`);
};

export const updateReportContent = async (reportContent: any) => {
  return getBackendSrv().put(
    `./api/plugins/msupply-datasource/resources/report-content/${reportContent?.id}`,
    reportContent
  );
};

export const getDatasources = async () => {
  return getBackendSrv().get(`./api/datasources`);
};

export const getSettings = async (id: any) => {
  return getBackendSrv().get(`/api/plugins/${id}/settings`);
};
