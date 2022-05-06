import React, { useEffect } from 'react';
import { SelectableValue, DateTime } from '@grafana/data';
import {
  Button,
  Field,
  FieldSet,
  FormAPI,
  Icon,
  InlineField,
  InlineFieldRow,
  Input,
  Select,
  TimeOfDayPicker,
} from '@grafana/ui';
import { getIntervals } from '../../constants';
import { Panel, ReportGroupType, ScheduleType } from 'types';
import { formatTimeToDate } from 'utils';
import { PanelList } from 'components';
import { Controller } from 'react-hook-form';

export const CreateScheduleFormPartial = ({
  register,
  errors,
  control,
  setValue,
  watch,
  reportGroups,
  isEditMode,
  defaultSchedule,
}: FormAPI<ScheduleType> & {
  reportGroups: ReportGroupType[] | undefined;
  isEditMode: boolean;
  defaultSchedule: ScheduleType;
}) => {
  useEffect(() => {
    // register('reportGroupID', { required: 'Report group is required' });
    // register('time', { required: 'time of day is required' });
    // register('interval', { required: 'Interval is required' });
  }, [register]);

  const getReportGroupOptions = (reportGroups: ReportGroupType[] | undefined) =>
    reportGroups?.map((reportGroup: ReportGroupType) => ({
      label: reportGroup.name,
      description: reportGroup.description,
      value: reportGroup,
    }));

  return (
    <>
      <FieldSet>
        <Field
          invalid={!!errors.name}
          error={errors.name && errors.name.message}
          label="Name"
          description="Name of the schedule"
        >
          <Input {...register('name', { required: 'Schedule name is required' })} id="schedule-name" />
        </Field>
        <Field label="description" description="Description of the schedule">
          <Input {...register('description')} id="schedule-description" />
        </Field>
      </FieldSet>

      <Field
        invalid={!!errors.reportGroupID}
        error={errors.reportGroupID && errors.reportGroupID.message}
        label="Report Group"
        description="Select a report group"
      >
        <Select
          value={reportGroups
            ?.filter((reportGroup: ReportGroupType) => reportGroup.id === watch('reportGroupID'))
            .map((reportGroup: ReportGroupType) => ({
              label: reportGroup.name,
              description: reportGroup.description,
              value: reportGroup,
            }))}
          options={getReportGroupOptions(reportGroups)}
          onChange={(selected: SelectableValue<ReportGroupType>) => {
            setValue('reportGroupID', selected?.value?.id ?? '');
          }}
          prefix={<Icon name="arrow-down" />}
        />
      </Field>

      <InlineFieldRow label={`Schedule time`} style={{ marginBottom: '30px' }}>
        <InlineField
          invalid={!!errors.interval}
          error={errors.interval && errors.interval.message}
          label="Interval"
          grow
          tooltip="Interval to queue the schedule on"
        >
          <Select
            value={getIntervals().filter((interval: any) => interval.value === watch('interval'))}
            options={getIntervals()}
            prefix={<Icon name="arrow-down" />}
            onChange={(option: any) => {
              setValue('interval', option.value);
            }}
          />
        </InlineField>
        <InlineField
          invalid={!!errors.time}
          error={errors.time && errors.time.message}
          label="Time of day"
          tooltip="Time of day to queue the schedule on"
          required
        >
          <TimeOfDayPicker
            value={formatTimeToDate(watch('time'))}
            onChange={(selected: DateTime) => {
              setValue('time', selected.format('HH:mm'));
            }}
          />
        </InlineField>
        {(watch('interval') || 0) > 2 && (
          <Field label="Report Day" description="The day to send the report in the month, half-year or year.">
            <Input type="number" {...register('day', { valueAsNumber: true })} id="schedule-day" width={40} />
          </Field>
        )}
      </InlineFieldRow>

      <Controller
        render={({ field: { onChange, value: selectedPanels } }) => {
          return (
            <PanelList
              panelListError={errors.panels}
              checkedPanels={selectedPanels}
              onPanelChecked={(panel: Panel) => {
                const updatedSelectedPanels = selectedPanels.includes(panel.id)
                  ? selectedPanels.filter((el: Number) => el !== panel.id)
                  : [...selectedPanels, panel.id];

                onChange(updatedSelectedPanels);
              }}
            />
          );
        }}
        name="panels"
        control={control}
        defaultValue={defaultSchedule.panelDetails.map((panelDetail) => panelDetail.panelID)}
      />

      <div className="gf-form-button-row">
        <Button type="submit" variant="primary">
          {isEditMode ? 'Update' : 'Create'} schedule
        </Button>
      </div>
    </>
  );
};
