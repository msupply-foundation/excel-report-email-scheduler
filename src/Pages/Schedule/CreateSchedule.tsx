import React, { useEffect } from 'react';
import { Button, Field, FieldSet, Form, Icon, Input, Select, TimeOfDayPicker } from '@grafana/ui';
import { Controller } from 'react-hook-form';

import { NAVIGATION_TITLE, NAVIGATION_SUBTITLE, ROUTES, getIntervals } from '../../constants';
import { Page, PanelList } from '../../components';
import { prefixRoute } from '../../utils';
import { Panel, ScheduleType } from 'types';
import { DateTime } from '@grafana/data';
import { useDatasourceID } from 'hooks';
import { useQuery } from 'react-query';
import { getPanels } from 'api/getPanels.api';

const defaultFormValues: ScheduleType = {
  id: '',
  name: '',
  description: '',
  interval: 0,
  timeOfDay: '',
  reportGroupID: '',
  day: 1,
  panels: [],
};

const CreateSchedule: React.FC = ({ history, match }: any) => {
  const datasourceID = useDatasourceID();

  const { data: panels } = useQuery<Panel[], Error>(['panels'], () => getPanels(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  useEffect(() => {
    console.log('datasourceID', datasourceID);
    console.log('panels', panels);
  }, [panels, datasourceID]);

  const submitCreateSchedule = () => {};

  const [defaultSchedule] = React.useState<ScheduleType>(defaultFormValues);

  return (
    <Page
      headerProps={{
        title: NAVIGATION_TITLE,
        subTitle: NAVIGATION_SUBTITLE,
        backButton: {
          icon: 'arrow-left',
          href: prefixRoute(ROUTES.SCHEDULERS),
        },
      }}
    >
      <Page.Contents>
        <Form
          onSubmit={submitCreateSchedule}
          validateOnMount={false}
          defaultValues={defaultSchedule}
          validateOn="onSubmit"
        >
          {({ register, errors, control, setValue, watch }) => {
            return (
              <>
                <FieldSet label={`New Schedule`}>
                  <Field
                    invalid={!!errors.name}
                    error={errors.name && errors.name.message}
                    label="Name"
                    description="Name of the schedule"
                  >
                    <Input
                      {...register('name', { required: 'Schedule name is required' })}
                      id="schedule-name"
                      width={60}
                    />
                  </Field>
                  <Field label="description" description="Description of the schedule">
                    <Input {...register('description')} id="schedule-description" width={60} />
                  </Field>
                </FieldSet>

                <FieldSet label={`Schedule time`}>
                  <Field label="Interval" description="Interval to queue the schedule on">
                    <Select
                      options={getIntervals()}
                      prefix={<Icon name="arrow-down" />}
                      onChange={(option: any) => {
                        setValue('interval', option.value);
                      }}
                    />
                  </Field>
                  <Field label="Time of day" description="Time of day to queue the schedule on">
                    <TimeOfDayPicker
                      onChange={(selected: DateTime) => {
                        setValue('timeOfDay', selected.format('HH:mm'));
                      }}
                    />
                  </Field>
                  {(watch('interval') || 0) > 2 && (
                    <Field label="Report Day" description="The day to send the report in the month, half-year or year.">
                      <Input type="number" {...register('day')} id="schedule-day" width={40} />
                    </Field>
                  )}
                </FieldSet>
                {panels && (
                  <Controller
                    render={({ field: { onChange, value: selectedPanels } }) => (
                      <PanelList
                        panels={panels}
                        panelListError={errors.panels}
                        checkedPanels={selectedPanels}
                        onPanelChecked={(event, panelID) => {
                          const updatedSelectedPanels = selectedPanels.includes(panelID)
                            ? selectedPanels.filter((el: any) => el !== panelID)
                            : [...selectedPanels, panelID];
                          onChange(updatedSelectedPanels);
                        }}
                      />
                    )}
                    name="panels"
                    control={control}
                  />
                )}

                <div className="gf-form-button-row">
                  <Button type="submit" variant="primary">
                    Create schedule
                  </Button>
                </div>
              </>
            );
          }}
        </Form>
      </Page.Contents>
    </Page>
  );
};

export { CreateSchedule };
