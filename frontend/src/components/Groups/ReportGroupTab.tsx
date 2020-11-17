import React, { FC, useState } from 'react';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { EditReportGroupModal } from './EditReportGroupModal';
import { createReportGroup } from 'api';
import { Button, Legend } from '@grafana/ui';
import { ReportGroupList } from './ReportGroupList';

import { AppRootProps } from '@grafana/data';
import { useOptimisticMutation } from 'hooks';
import { AppData, CreateReportGroupVariables, ReportGroup } from 'common/types';
import { AppDataContext } from '../../containers';
import { queryCache } from 'react-query';

interface Props extends AppRootProps {}

const adjustButtonToRight = css`
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
`;

export const ReportGroupTab: FC<Props> = ({ meta }) => {
  const appData = meta?.jsonData as AppData;
  const [activeGroup, setActiveGroup] = useState<string | null | undefined>(null);
  const groups = queryCache.getQueryData<ReportGroup[]>('reportGroup');
  const group = groups?.find((group: ReportGroup) => group.id === activeGroup);

  const [newReportGroup] = useOptimisticMutation<ReportGroup[], ReportGroup, CreateReportGroupVariables, ReportGroup[]>(
    ['reportGroup'],
    createReportGroup,
    () => ({ id: '', name: intl.get('new_report_group'), description: '' }),
    (prevState, optimisticValue) => {
      if (prevState) {
        return [...prevState, optimisticValue];
      }
      return [optimisticValue];
    },
    []
  );

  return (
    <div>
      <AppDataContext.Provider value={appData}>
        <div className={adjustButtonToRight}>
          <Legend>Report Groups</Legend>
          <Button onClick={newReportGroup} variant="primary">
            {intl.get('add_report_group')}
          </Button>
        </div>
        <ReportGroupList onRowPress={setActiveGroup} />
        {group && groups && (
          <EditReportGroupModal
            datasourceID={meta?.jsonData?.datasourceID}
            reportGroup={group}
            isOpen={!!activeGroup}
            onClose={() => setActiveGroup(null)}
          />
        )}
      </AppDataContext.Provider>
    </div>
  );
};
