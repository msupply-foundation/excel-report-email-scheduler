import React, { FC, useState } from 'react';
import { queryCache, useMutation } from 'react-query';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { EditReportGroupModal } from './EditReportGroupModal';
import { createReportGroup } from 'api';
import { Button, Legend } from '@grafana/ui';
import { ReportGroupList } from './ReportGroupList';
import { ReportGroup } from './ReportSchedulesTab';
import { AppRootProps } from '@grafana/data';

interface Props extends AppRootProps {}

const adjustButtonToRight = css`
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
`;

export const ReportGroupTab: FC<Props> = ({ meta }) => {
  const [activeGroup, setActiveGroup] = useState<ReportGroup | null>(null);
  const [newReportGroup] = useMutation(createReportGroup, {
    onSuccess: () => queryCache.refetchQueries(['reportGroup']),
  });

  return (
    <div>
      <div className={adjustButtonToRight}>
        <Legend>Report Groups</Legend>
        <Button onClick={newReportGroup} variant="primary">
          {intl.get('add_report_group')}
        </Button>
      </div>
      <ReportGroupList onRowPress={setActiveGroup} />
      {activeGroup && (
        <EditReportGroupModal
          datasourceID={meta?.jsonData?.datasourceID}
          reportGroup={activeGroup}
          isOpen={!!activeGroup}
          onClose={() => setActiveGroup(null)}
        />
      )}
    </div>
  );
};
