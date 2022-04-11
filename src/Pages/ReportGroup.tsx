import React from 'react';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { Button, useStyles2 } from '@grafana/ui';
import intl from 'react-intl-universal';
import { EmptyListCTA } from 'components/common';
import { prefixRoute } from '../utils';
import { ROUTES } from '../constants';

const EmptyList = () => {
  return (
    <EmptyListCTA
      title="You haven't created any report groups yet."
      buttonTitle={'Create new report group'}
      buttonIcon="users-alt"
      buttonLink={`${prefixRoute(ROUTES.REPORT_GROUP)}/create`}
      proTip="Report groups are groups containing users to send the report to."
      proTipLink=""
      proTipLinkTitle=""
      proTipTarget="_blank"
    />
  );
};

const ReportGroup = () => {
  const styles = useStyles2(getStyles);

  return <EmptyList />;

  return (
    <div>
      <div className={styles.adjustButtonToRight}>
        <Button onClick={() => {}} variant="primary">
          {intl.get('add_report_group')}
        </Button>
      </div>
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  marginTop: css`
    margin-top: ${theme.spacing(2)};
  `,
  adjustButtonToRight: css`
    display: flex;
    justify-content: flex-end;
    margin-bottom: 10px;
  `,
});

export { ReportGroup, EmptyList };
