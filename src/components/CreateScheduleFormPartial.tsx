import React, { useEffect } from 'react';
import { SelectableValue, DateTime } from '@grafana/data';
import { Button, Field, FieldSet, FormAPI, Icon, Input, Select, TimeOfDayPicker } from '@grafana/ui';
import { getIntervals } from '../constants';
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

  return (
    <>
      <FieldSet label={`${isEditMode ? 'Edit "' + defaultSchedule?.name + '"' : 'New'} Schedule`}>
        <Field
          invalid={!!errors.name}
          error={errors.name && errors.name.message}
          label="Name"
          description="Name of the schedule"
        >
          <Input {...register('name', { required: 'Schedule name is required' })} id="schedule-name" width={60} />
        </Field>
        <Field label="description" description="Description of the schedule">
          <Input {...register('description')} id="schedule-description" width={60} />
        </Field>
      </FieldSet>

      <Field
        invalid={!!errors.reportGroupID}
        error={errors.reportGroupID && errors.reportGroupID.message}
        label="Report Group"
        description="Select a report group"
      >
        <Select
          options={reportGroups?.map((reportGroup: ReportGroupType) => ({
            label: reportGroup.name,
            description: reportGroup.description,
            value: reportGroup,
          }))}
          onChange={(selected: SelectableValue<ReportGroupType>) => {
            setValue('reportGroupID', selected?.value?.id ?? '');
          }}
          prefix={<Icon name="arrow-down" />}
        />
      </Field>

      <FieldSet label={`Schedule time`}>
        <Field
          invalid={!!errors.interval}
          error={errors.interval && errors.interval.message}
          label="Interval"
          description="Interval to queue the schedule on"
        >
          <Select
            options={getIntervals()}
            prefix={<Icon name="arrow-down" />}
            onChange={(option: any) => {
              setValue('interval', option.value);
            }}
          />
        </Field>
        <Field
          invalid={!!errors.time}
          error={errors.time && errors.time.message}
          label="Time of day"
          description="Time of day to queue the schedule on"
        >
          <TimeOfDayPicker
            value={formatTimeToDate(watch('time'))}
            onChange={(selected: DateTime) => {
              setValue('time', selected.format('HH:mm'));
            }}
          />
        </Field>
        {/* {(watch('interval') || 0) > 2 && ( */}
        <Field label="Report Day" description="The day to send the report in the month, half-year or year.">
          <Input type="number" {...register('day', { valueAsNumber: true })} id="schedule-day" width={40} />
        </Field>
        {/* )} */}
      </FieldSet>

      <Controller
        render={({ field: { onChange, value: selectedPanels } }) => {
          console.log('selectedPanels', selectedPanels);
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
      />

      <div className="gf-form-button-row">
        <Button type="submit" variant="primary">
          Create schedule
        </Button>
      </div>
    </>
  );
};
