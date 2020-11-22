import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Icon, InlineFormLabel, Input, Legend, Select, Tooltip } from '@grafana/ui';

import { ScheduleKey } from 'common/enums';
import { ReportGroup, Schedule } from 'common/types';
import { useQuery } from 'react-query';
import { getReportGroups } from 'api';
import { SelectableValue } from '@grafana/data';

const getIntervals = () => [
  { label: intl.get('daily'), value: 60 * 60 * 24 },
  { label: intl.get('weekly'), value: 60 * 60 * 24 * 7 },
  { label: intl.get('fortnightly'), value: 60 * 60 * 24 * 14 },
  { label: intl.get('monthly'), value: 60 * 60 * 24 * 30 },
  { label: intl.get('quarterly'), value: 60 * 60 * 24 * 30 * 6 },
  { label: intl.get('yearly'), value: 60 * 60 * 24 * 30 * 12 },
];

const container = css`
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  flex: 10;
  padding-right: 30px;
  flex-direction: column;
`;

const flexWrapping = css`
  display: flex;
  flex-direction: row;
  flex: 1;
`;

const flexContainer = css`
  display: flex;
  flex: 1;
  min-width: 50%;
`;

type Props = {
  onUpdate: (key: ScheduleKey, value: string) => void;
  schedule: Schedule;
};

export const EditScheduleForm: FC<Props> = ({ onUpdate, schedule }) => {
  const { data: reportGroups } = useQuery('reportGroup', getReportGroups);

  return (
    <div className={container}>
      <div style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
        <Tooltip placement="top" content={intl.get('edit_details_schedule_tooltip')} theme={'info'}>
          <Icon
            name="info-circle"
            size="sm"
            style={{ marginLeft: '10px', marginRight: '10px', marginBottom: '16px' }}
          />
        </Tooltip>
        <Legend>{intl.get('edit_details')}</Legend>
      </div>

      <div className={flexWrapping}>
        <InlineFormLabel tooltip={intl.get('group_name')}>{intl.get('name')}</InlineFormLabel>

        <div className={flexContainer}>
          <Input
            onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.NAME, value)}
            name={intl.get('name')}
            defaultValue={schedule.name}
            placeholder={intl.get('name')}
            css=""
          />
        </div>
      </div>

      <div className={flexWrapping}>
        <InlineFormLabel tooltip={intl.get('group_description')}>{intl.get('description')}</InlineFormLabel>

        <div className={flexContainer}>
          <Input
            onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.DESCRIPTION, value)}
            name={intl.get('description')}
            defaultValue={schedule.description}
            placeholder={intl.get('description')}
            css=""
          />
        </div>
      </div>

      <div className={flexWrapping}>
        <InlineFormLabel tooltip={intl.get('report_interval_description')}>
          {intl.get('report_interval')}
        </InlineFormLabel>

        <div className={flexContainer}>
          <Select
            value={getIntervals().filter((interval: any) => interval.value === schedule?.interval)}
            options={getIntervals()}
            onChange={(selected: SelectableValue) => {
              onUpdate(ScheduleKey.INTERVAL, selected.value);
            }}
          />
        </div>
      </div>

      <div className={flexWrapping}>
        <InlineFormLabel tooltip={intl.get('report_group_description')}>{intl.get('report_group')}</InlineFormLabel>

        <div className={flexContainer}>
          <Select
            value={reportGroups
              ?.filter((reportGroup: ReportGroup) => reportGroup.id === schedule.reportGroupID)
              .map((reportGroup: ReportGroup) => ({
                label: reportGroup.name,
                description: reportGroup.description,
                value: reportGroup,
              }))}
            options={reportGroups?.map((reportGroup: ReportGroup) => ({
              label: reportGroup.name,
              description: reportGroup.description,
              value: reportGroup,
            }))}
            onChange={(selected: SelectableValue<ReportGroup>) => {
              onUpdate(ScheduleKey.REPORT_GROUP_ID, selected?.value?.id ?? '');
            }}
          />
        </div>
      </div>
    </div>
  );
};
