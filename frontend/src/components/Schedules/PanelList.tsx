import React, { FC, useCallback } from 'react';
import intl from 'react-intl-universal';
import classNames from 'classnames';

import { createReportContent, deleteReportContent, getPanels, getReportContent } from '../../api';

import { useQuery } from 'react-query';
import { Schedule, ReportContent, Panel, CreateContentVars } from 'common/types';

import { useOptimisticMutation } from 'hooks/useOptimisticMutation';
import { PanelItem } from './PanelItem';
import { Icon, Legend, Tooltip } from '@grafana/ui';
import { useDatasourceID } from 'hooks';
import { useWindowSize } from 'hooks/useWindowResize';

const listStyle = classNames({
  'card-section': true,
  'card-list-layout-grid': true,
  'card-list-layout-list': true,
  'card-list': true,
});

type Props = {
  schedule: Schedule;
};

const findMatchingContent = (reportContents: ReportContent[], panel: Panel) =>
  reportContents?.find((reportContent: ReportContent) => {
    const { panelID, dashboardID } = reportContent;
    const { id, dashboardID: panelDashboardID } = panel;
    return id === panelID && dashboardID === panelDashboardID;
  }) ?? null;

export const PanelList: FC<Props> = ({ schedule }) => {
  const { height } = useWindowSize();
  const datasourceID = useDatasourceID();
  const { id: scheduleID } = schedule;

  const { data: panels } = useQuery<Panel[], Error>(['panels'], () => getPanels(datasourceID));

  const { data: reportContents } = useQuery<ReportContent[], Error>(['reportContent', scheduleID], getReportContent);

  const getMatchingContent = useCallback(
    (panel: Panel) => {
      if (reportContents) {
        return findMatchingContent(reportContents, panel);
      } else {
        return null;
      }
    },
    [reportContents]
  );

  const [createContent] = useOptimisticMutation<ReportContent[], ReportContent, CreateContentVars, ReportContent[]>(
    ['reportContent', scheduleID],
    createReportContent,
    (variables: CreateContentVars): ReportContent => ({ ...variables, lookback: 0, id: '' }),
    (prevData: ReportContent[] | undefined, optimistic: ReportContent): ReportContent[] => {
      if (prevData) {
        return [...prevData, optimistic];
      }
      return [optimistic];
    },
    []
  );

  const [deleteContent] = useOptimisticMutation<ReportContent[], ReportContent, ReportContent, ReportContent[]>(
    ['reportContent', scheduleID],
    deleteReportContent,
    (reportContents: ReportContent): ReportContent => reportContents,
    (prevState: ReportContent[] | undefined, optimisticValue: ReportContent) => {
      if (prevState) {
        return prevState.filter(({ id }) => id !== optimisticValue.id);
      }
      return prevState;
    },
    []
  );

  const onTogglePanel = async (panel: Panel) => {
    const matchingContent = getMatchingContent(panel);
    // Checking an ID exists checks if this is an optimistic
    // created content or a real report content.
    if (matchingContent?.id) {
      await deleteContent(matchingContent);
    } else {
      await createContent({
        scheduleID: scheduleID ?? '',
        panelID: panel?.id,
        dashboardID: panel?.dashboardID,
        variables: '',
      });
    }
  };

  return (
    <div style={{ marginTop: '25px' }}>
      <div style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
        <Tooltip placement="top" content={intl.get('available_panels_tooltip')} theme={'info'}>
          <Icon
            name="info-circle"
            size="sm"
            style={{ marginLeft: '10px', marginRight: '10px', marginBottom: '16px' }}
          />
        </Tooltip>
        <Legend>{intl.get('available_panels')}</Legend>
      </div>
      <ol className={listStyle} style={{ maxHeight: `${(height ?? 0) / 2}px`, overflow: 'scroll' }}>
        {panels?.map((panel: any) => {
          const matchingContent = getMatchingContent(panel);

          return (
            <PanelItem
              reportContent={matchingContent}
              panel={panel}
              scheduleID={scheduleID}
              onToggle={onTogglePanel}
              key={scheduleID}
            />
          );
        })}
      </ol>
    </div>
  );
};
