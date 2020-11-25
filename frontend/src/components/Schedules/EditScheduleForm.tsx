import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Icon, InlineFormLabel, Input, Legend, Select, Tooltip, TimeOfDayPicker } from '@grafana/ui';

import { ScheduleKey } from 'common/enums';
import { ReportGroup, Schedule } from 'common/types';
import { useQuery } from 'react-query';
import { getReportGroups } from 'api';
import { DateTime, dateTime, SelectableValue } from '@grafana/data';

const getIntervals = () => [
  { label: intl.get('daily'), value: 0 },
  { label: intl.get('weekly'), value: 1 },
  { label: intl.get('fortnightly'), value: 2 },
  { label: intl.get('monthly'), value: 3 },
  { label: intl.get('quarterly'), value: 4 },
  { label: intl.get('yearly'), value: 5 },
];

const formatTimeToDate = (time?: string) => {
  const now: DateTime = dateTime(Date.now());
  const d: DateTime = dateTime(now.format('YYYY-MM-DD') + ' ' + time, 'YYYY-MM-DD HH:mm');
  return d.isValid() ? d : undefined;
};

const container = css`
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  flex: 1;
  padding-right: 30px;
  flex-direction: row;
`;

const flexWrapping = css`
  display: flex;
  flex-direction: row;
  flex: 1;
  flex-basis: calc(50% - 20px);
  margin-left: 20px;
`;

type Props = {
  onUpdate: (key: ScheduleKey, value: string | number) => void;
  schedule: Schedule;
};

export const EditScheduleForm: FC<Props> = ({ onUpdate, schedule }) => {
  const { data: reportGroups } = useQuery('reportGroup', getReportGroups);

  return (
    <div style={{ width: '100%' }}>
      <div>
        <Tooltip placement="top" content={intl.get('edit_details_schedule_tooltip')} theme={'info'}>
          <Icon
            name="info-circle"
            size="sm"
            style={{ marginLeft: '10px', marginRight: '10px', marginBottom: '16px' }}
          />
        </Tooltip>
        <Legend>{intl.get('edit_details')}</Legend>
      </div>
      <div className={container}>
        <div className={flexWrapping}>
          <InlineFormLabel tooltip={intl.get('group_name')}>{intl.get('name')}</InlineFormLabel>

          <Input
            onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.NAME, value)}
            name={intl.get('name')}
            defaultValue={schedule.name}
            placeholder={intl.get('name')}
            css=""
          />
        </div>

        <div className={flexWrapping}>
          <InlineFormLabel tooltip={intl.get('report_interval_description')}>
            {intl.get('report_interval')}
          </InlineFormLabel>

          <Select
            value={getIntervals().filter((interval: any) => interval.value === schedule?.interval)}
            options={getIntervals()}
            onChange={(selected: SelectableValue) => {
              onUpdate(ScheduleKey.INTERVAL, selected.value);
            }}
          />
        </div>

        <div className={flexWrapping}>
          <InlineFormLabel tooltip={intl.get('group_description')}>{intl.get('description')}</InlineFormLabel>

          <Input
            onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.DESCRIPTION, value)}
            name={intl.get('description')}
            defaultValue={schedule.description}
            placeholder={intl.get('description')}
            css=""
          />
        </div>

        <div className={flexWrapping}>
          <InlineFormLabel tooltip={intl.get('report_time_description')}>{intl.get('report_time')}</InlineFormLabel>

          <TimeOfDayPicker
            onChange={(selected: DateTime) => {
              onUpdate(ScheduleKey.TIME_OF_DAY, selected.format('HH:mm'));
            }}
            value={formatTimeToDate(schedule?.time)}
          />
        </div>

        <div className={flexWrapping}>
          <InlineFormLabel tooltip={intl.get('report_group_description')}>{intl.get('report_group')}</InlineFormLabel>

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

        {(schedule.interval || 0) > 2 ? (
          <div className={flexWrapping}>
            <InlineFormLabel tooltip={intl.get('report_day_description')}>{intl.get('report_day')}</InlineFormLabel>

            <Input
              css=""
              defaultValue={1}
              min={1}
              name={intl.get('report_day')}
              onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.DAY_OF_INTERVAL, parseInt(value, 10))}
              placeholder={intl.get('report_day')}
              type="number"
              value={schedule?.day || 1}
            />
          </div>
        ) : (
          <div className={flexWrapping}></div>
        )}
      </div>
    </div>
  );
};
