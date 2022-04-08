import * as React from 'react';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { Button, useStyles2 } from '@grafana/ui';
import intl from 'react-intl-universal';

const ReportGroup = () => {
  const styles = useStyles2(getStyles);

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

export { ReportGroup };
