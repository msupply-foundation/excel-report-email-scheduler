import { Form, FormAPI } from '@grafana/ui';
import { createSchedule, getReportGroups } from 'api';
import { CreateScheduleFormPartial } from 'components';
import { PLUGIN_BASE_URL } from '../../constants';
import { PanelContext } from 'context';
import React, { useContext } from 'react';
import { useMutation, useQuery } from 'react-query';
import { ScheduleType, ReportGroupType, PanelDetails } from 'types';

const defaultFormValues: ScheduleType = {
  id: '',
  name: '',
  description: '',
  interval: 0,
  time: '',
  reportGroupID: '',
  day: 1,
  panels: [],
  panelDetails: [],
};

const CreateScheduleForm: React.FC = ({ history, match }: any) => {
  const { panelDetails } = useContext(PanelContext);

  const { data: reportGroups } = useQuery<ReportGroupType[], Error>(`reportGroups`, getReportGroups, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const createScheduleMutation = useMutation((newSchedule: ScheduleType) => createSchedule(newSchedule), {
    onSuccess: () => {
      history.push(`${PLUGIN_BASE_URL}/schedulers/`);
      return;
    },
  });

  const submitCreateSchedule = (data: ScheduleType) => {
    console.log(data);
    const selectedPanels = panelDetails.filter((detail: PanelDetails) => data.panels.includes(detail.panelID));
    data.panelDetails = selectedPanels;

    console.log(data);
    createScheduleMutation.mutate(data);
  };

  return (
    <Form
      onSubmit={submitCreateSchedule}
      validateOnMount={false}
      defaultValues={defaultFormValues}
      validateOn="onSubmit"
    >
      {(props: FormAPI<ScheduleType>) => {
        return <CreateScheduleFormPartial {...props} reportGroups={reportGroups} />;
      }}
    </Form>
  );
};

export { CreateScheduleForm };
