import React from 'react';

import { NAVIGATION_TITLE, NAVIGATION_SUBTITLE, ROUTES } from '../../constants';
import { Page } from '../../components';
import { prefixRoute } from '../../utils';
import { PanelProvider } from 'context';
import { CreateScheduleForm } from 'components/schedule/CreateScheduleForm';
import { RouteComponentProps } from 'react-router-dom';
import { PageToolbar } from '@grafana/ui';

const CreateSchedule: React.FC<RouteComponentProps> = () => (
  <Page
    headerProps={{
      title: NAVIGATION_TITLE,
      subTitle: NAVIGATION_SUBTITLE,
    }}
  >
    <PanelProvider>
      <Page.Contents>
        <CreateScheduleForm />
      </Page.Contents>
    </PanelProvider>
  </Page>
);

export { CreateSchedule };
