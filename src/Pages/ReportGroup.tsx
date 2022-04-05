import * as React from 'react';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { LinkButton, useStyles2 } from '@grafana/ui';
import { ROUTES } from '../constants';
import { prefixRoute } from '../utils/navigation';

const ReportGroup = () => {
  const s = useStyles2(getStyles);

  return (
    <div>
      This is page one.
      <div className={s.marginTop}>
        <LinkButton href={prefixRoute(ROUTES.REPORT_GROUP)}>Full-width page example</LinkButton>
      </div>
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  marginTop: css`
    margin-top: ${theme.spacing(2)};
  `,
});

export { ReportGroup };
