import React from 'react';
import { Form } from '@grafana/ui';

import { NAVIGATION_TITLE, NAVIGATION_SUBTITLE, ROUTES } from '../../constants';
import { CreateScheduleFormPartial, Page } from '../../components';
import { prefixRoute } from '../../utils';
import { Panel, ReportGroupType, ScheduleType } from 'types';
import { useDatasourceID } from 'hooks';
import { useQuery } from 'react-query';
import { getPanels } from 'api/getPanels.api';
import { getReportGroups } from 'api';

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

  const { data: reportGroups } = useQuery<ReportGroupType[], Error>(`reportGroups`, getReportGroups, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const submitCreateSchedule = (data: ScheduleType) => {
    console.log(data);
  };

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
          {(props) => {
            return <CreateScheduleFormPartial {...props} panels={panels} reportGroups={reportGroups} />;
          }}
        </Form>
      </Page.Contents>
    </Page>
  );
};

export { CreateSchedule };
