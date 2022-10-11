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
import { getIntervals, getWeekDays, getDateFormat, getDatePosition } from '../../constants';
import { Panel, PanelListSelectedType, ReportGroupType, ScheduleType } from 'types';
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
    register('reportGroupID', { required: 'Report group is required' });
    register('time', { required: 'time of day is required' });
    register('interval', { required: 'Interval is required' });
    register('dateFormat');
    register('datePosition');
  }, [register]);

  const getReportGroupOptions = (reportGroups: ReportGroupType[] | undefined) =>
    reportGroups?.map((reportGroup: ReportGroupType) => ({
      label: reportGroup.name,
      description: reportGroup.description,
      value: reportGroup,
    }));

  const renderInterval = () => {
    switch (watch('interval')) {
      case 0:
        return false;
      case 1:
        return (
          <InlineFieldRow label="Report Day">
            <InlineField
            label="Report Day"
            grow
            tooltip="The day to send the report in Weekly"
            >
            <Select
              value={getWeekDays().filter((inter: any) => inter.value === watch('day'))}
              options={getWeekDays()}
              prefix={<Icon name="arrow-down" />}
              onChange={(option: any) => {
                setValue('day', option.value);
              }}
            />
            </InlineField>
          </InlineFieldRow>
        );
      case 2:
        return (
          <InlineFieldRow label="Report Day">
          <InlineField
          label="Report Day"
          grow
          tooltip="The day to send the report in Fortnightly"
          >
            <Select
              value={watch('day')}
              options={[...Array(14).keys()].map((key) => ({ label: (key + 1).toString(), value: key + 1 }))}
              prefix={<Icon name="arrow-down" />}
              onChange={(option: any) => {
                setValue('day', option.value);
              }}
            />
          </InlineField>
        </InlineFieldRow>
        );
      default:
        return (
          <InlineFieldRow label="Report Day">
          <InlineField
          label="Report Day"
          grow
          tooltip="The day to send the report in the month, half-year or year."
          >
            <Input type="number" {...register('day', { valueAsNumber: true })} id="schedule-day" width={40} />
          </InlineField>
        </InlineFieldRow>
        );
    }
  };

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
        {renderInterval()}
        <InlineField
          label="Date Format"
          grow
          tooltip="Date format to attach in filename."
        >
          <Select
            value={getDateFormat().filter((dateFormat: any) => {
              return dateFormat.value === watch('dateFormat');
            })}
            options={getDateFormat()}
            prefix={<Icon name="arrow-down" />}
            onChange={(option: any) => {
              setValue('dateFormat', option.value);
            }}
          />
        </InlineField>
        <InlineField
          label="Add Date To"
          grow
          tooltip="Add date either in start or end of filename."
        >
          <Select
            value={getDatePosition().filter((datePosition: any) => {
              return datePosition.value === watch('datePosition');
            })}
            options={getDatePosition()}
            prefix={<Icon name="arrow-down" />}
            onChange={(option: any) => {
              setValue('datePosition', option.value);
            }}
          />
        </InlineField>
      </InlineFieldRow>
      <Controller
        render={({ field: { onChange, value: selectedPanels } }) => {
          return (
            <PanelList
              panelListError={errors.panels}
              checkedPanels={selectedPanels}
              onPanelChecked={(panel: Panel) => {
                const foundElement = selectedPanels.find(
                  (selectedPanel: PanelListSelectedType) =>
                    selectedPanel.panelID === panel.id && selectedPanel.dashboardID === panel.dashboardID
                );

                const updatedSelectedPanels = foundElement
                  ? selectedPanels.filter(
                      (selectedPanel: PanelListSelectedType) =>
                        !(selectedPanel.panelID === panel.id && selectedPanel.dashboardID === panel.dashboardID)
                    )
                  : [...selectedPanels, { panelID: panel.id, dashboardID: panel.dashboardID }];

                onChange(updatedSelectedPanels);
              }}
            />
          );
        }}
        name="panels"
        control={control}
        defaultValue={defaultSchedule.panelDetails.map((panelDetail) => ({
          panelID: panelDetail.panelID,
          dashboardID: panelDetail.dashboardID,
        }))}
      />

      <div className="gf-form-button-row">
        <Button type="submit" variant="primary">
          {isEditMode ? 'Update' : 'Create'} schedule
        </Button>
      </div>
    </>
  );
};
