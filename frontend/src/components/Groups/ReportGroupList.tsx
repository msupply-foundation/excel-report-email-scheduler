import React, { FC } from 'react';
import classNames from 'classnames';
import { Spinner } from '@grafana/ui';
import { getReportGroups } from '../../api';
import { useQuery } from 'react-query';
import { ReportGroup } from '../Schedules/ReportSchedulesTab';

const listStyle = classNames({
  'card-section': true,
  'card-list-layout-grid': true,
  'card-list-layout-list': true,
  'card-list': true,
});

type Props = {
  onRowPress: (reportGroup: string | undefined) => void;
};

export const ReportGroupList: FC<Props> = ({ onRowPress }) => {
  const { data: reportGroups, isLoading } = useQuery('reportGroup', getReportGroups);

  return isLoading ? (
    <Spinner />
  ) : (
    <ol className={listStyle}>
      {reportGroups?.map((reportGroup: ReportGroup) => {
        const { name, description } = reportGroup;
        return (
          <li className="card-item-wrapper" style={{ cursor: 'pointer' }}>
            <div className="card-item" onClick={() => onRowPress(reportGroup?.id)}>
              <div className="card-item-name">{name}</div>
              {description && <div className="card-item-type">{description}</div>}
            </div>
          </li>
        );
      })}
    </ol>
  );
};
