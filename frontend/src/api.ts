import { panelUsesVariable } from './common/utils/checkers';
import { DashboardResponse, DashboardMeta, CreateContentVars } from './common/types';
import { getBackendSrv } from '@grafana/runtime';
import { Variable, Panel, ReportGroupMember, Schedule, Store } from 'common/types';

export const getRecipients = () => getBackendSrv().get('api/plugins/msupply-datasource/resources/report-recipient');

export const getGroupAssignments = (key: string, groupId: string) =>
  getBackendSrv().get(`api/plugins/msupply-datasource/resources/report-group-membership/?group-id=${groupId}`);

export const getGroupMembers = (key: string, groupId: string) =>
  getBackendSrv().get(`api/plugins/msupply-datasource/resources/report-group-membership/?group-id=${groupId}`);

export const getReportGroups = () => getBackendSrv().get('api/plugins/msupply-datasource/resources/report-group');

export const getUsers = (datasourceID: number) => {
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

export const deleteReportGroup = async (reportGroup: any) => {
  return getBackendSrv().delete(`./api/plugins/msupply-datasource/resources/report-group/${reportGroup?.id}`);
};

export const createReportGroup = async (ReportGroup: any) => {
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
  const content: any[] = await getBackendSrv().get(
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

export const getPanels = async (): Promise<Panel[]> => {
  const dashboards = (await getDashboards()) ?? [];
  const panels: Panel[] = dashboards
    .filter(({ panels }) => panels?.length > 0)
    .map<Panel[]>(({ panels, templating, uid }) => {
      const mappedPanels = panels
        .filter(({ type }) => type === 'table')
        .map(rawPanel => {
          const { targets, description, title, id, type } = rawPanel;

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

          const mappedPanel = { rawSql, description, title, id, variables, dashboardID, type };
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
