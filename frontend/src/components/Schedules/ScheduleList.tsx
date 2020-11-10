import React, { FC } from 'react';
import classNames from 'classnames';
import { Spinner } from '@grafana/ui';

import { getSchedules } from 'api';
import { useQuery } from 'react-query';
import { Schedule } from 'common/types';

type Props = {
  onRowPress: (toggle: any) => void;
};

const listStyle = classNames({
  'card-section': true,
  'card-list-layout-grid': true,
  'card-list-layout-list': true,
  'card-list': true,
});

export const ScheduleList: FC<Props> = ({ onRowPress }) => {
  const { data: schedules, isLoading } = useQuery('reportSchedules', getSchedules);

  return isLoading ? (
    <Spinner />
  ) : (
    <ol className={listStyle}>
      {schedules?.map((schedule: Schedule) => {
        const { name, description } = schedule;
        return (
          <li className="card-item-wrapper" style={{ cursor: 'pointer' }}>
            <div className={'card-item'} onClick={() => onRowPress(schedule)}>
              <div className="card-item-details">
                <div className="card-item-name">{name}</div>
                <div className="card-item-type">{description}</div>
              </div>
            </div>
          </li>
        );
      })}
    </ol>
  );
};
