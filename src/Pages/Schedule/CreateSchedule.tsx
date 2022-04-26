import React from 'react';
import { Form } from '@grafana/ui';

import { NAVIGATION_TITLE, NAVIGATION_SUBTITLE, ROUTES } from '../../constants';
import { Page } from '../../components';
import { prefixRoute } from '../../utils';

// const defaultFormValues: any = {
//   id: '',
//   name: '',
//   description: '',
//   members: [],
// };

const CreateSchedule: React.FC = ({ history, match }: any) => {
  const submitCreateSchedule = () => {};

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
        <Form onSubmit={submitCreateSchedule} validateOnMount={false} validateOn="onSubmit">
          {({ register, errors, control }) => {
            return <></>;
          }}
        </Form>
      </Page.Contents>
    </Page>
  );
};

export { CreateSchedule };
