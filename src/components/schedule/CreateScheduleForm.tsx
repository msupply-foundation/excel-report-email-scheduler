import { Form } from '@grafana/ui';
import { getReportGroups } from 'api';
import { CreateScheduleFormPartial } from 'components';
import { PanelContext } from 'context';
import React, { useContext } from 'react';
import { useQuery } from 'react-query';
import { ScheduleType, ReportGroupType, PanelDetails } from 'types';

const defaultFormValues: ScheduleType = {
  id: '',
  name: '',
  description: '',
  interval: 0,
  timeOfDay: '',
  reportGroupID: '',
  day: 1,
  panels: [],
  panelDetails: [],
};

const CreateScheduleForm: React.FC = () => {
  const { panelDetails } = useContext(PanelContext);

  const { data: reportGroups } = useQuery<ReportGroupType[], Error>(`reportGroups`, getReportGroups, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const submitCreateSchedule = (data: ScheduleType) => {
    console.log('panelDetails', panelDetails);
    const selectedPanels = panelDetails.filter((detail: PanelDetails) => data.panels.includes(detail.panelID));
    data.panelDetails = selectedPanels;
    console.log(data);
  };

  const [defaultSchedule] = React.useState<ScheduleType>(defaultFormValues);

  return (
    <Form onSubmit={submitCreateSchedule} validateOnMount={false} defaultValues={defaultSchedule} validateOn="onSubmit">
      {(props: any) => {
        return <CreateScheduleFormPartial {...props} reportGroups={reportGroups} />;
      }}
    </Form>
  );
};

export { CreateScheduleForm };
