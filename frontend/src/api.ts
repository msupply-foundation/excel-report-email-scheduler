import { getBackendSrv } from '@grafana/runtime';
import { ReportGroupMember, Schedule } from 'common/types';

export const getRecipients = () => getBackendSrv().get('api/plugins/msupply-datasource/resources/report-recipient');

export const getDashboards = () => getBackendSrv().get('/api/search');

export const getGroupAssignments = (key: string, groupId: string) =>
  getBackendSrv().get(`api/plugins/msupply-datasource/resources/report-group-membership/?group-id=${groupId}`);

export const getGroupMembers = (key: string, groupId: string) =>
  getBackendSrv().get(`api/plugins/msupply-datasource/resources/report-group-membership/?group-id=${groupId}`);

export const getReportGroups = () => getBackendSrv().get('api/plugins/msupply-datasource/resources/report-group');

export const getUsers = () => {
  return getBackendSrv()
    .post('/api/tsdb/query', {
      queries: [
        {
          datasourceId: 1,
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

export const getStores = () => {
  return getBackendSrv()
    .post('/api/tsdb/query', {
      queries: [
        {
          datasourceId: 1,
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

  const membership = { reportRecipientID: user.id, reportGroupID };
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
  return content.reduce((acc: any, value: any) => ({ ...acc, [value.panelID]: value }), {});
};

// TODO: Make this not terrible!
// Searches for all dashboards, then does a query for each dashboard and maps each
// dashboards panel if it is a table format panel.
export const getPanels = async () => {
  return (
    await Promise.all(
      (await getBackendSrv().get('./api/search')).map((dash: any) => {
        return getBackendSrv().get(`./api/dashboards/uid/${dash.uid}`);
      })
    )
  )
    .map((det: any) => det.dashboard)
    .map((dashboard: any) => dashboard?.panels?.filter((panel: any) => panel?.targets?.[0].format === 'table') ?? [])
    .flat();
};

export const updateSchedule = async (schedule: Schedule) => {
  return getBackendSrv().put(`./api/plugins/msupply-datasource/resources/schedule/${schedule?.id}`, schedule);
};

export const createReportContent = async (contents: any) => {
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
