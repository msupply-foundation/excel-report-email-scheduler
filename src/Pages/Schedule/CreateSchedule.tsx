import React from 'react';

import { NAVIGATION_TITLE, NAVIGATION_SUBTITLE, ROUTES } from '../../constants';
import { Page } from '../../components';
import { prefixRoute } from '../../utils';
import { PanelProvider } from 'context';
import { CreateScheduleForm } from 'components/schedule/CreateScheduleForm';
import { RouteComponentProps } from 'react-router-dom';

const CreateSchedule: React.FC<RouteComponentProps> = () => (
  <Page
    headerProps={{
      title: NAVIGATION_TITLE,
      subTitle: NAVIGATION_SUBTITLE,
      backButton: {
        icon: 'arrow-left',
        href: prefixRoute(ROUTES.SCHEDULES),
      },
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
