import * as React from 'react';
// import { css } from '@emotion/css';
// import { GrafanaTheme2 } from '@grafana/data';
// import { useStyles2 } from '@grafana/ui';
import { ROUTES } from '../../constants';
import { prefixRoute } from '../../utils/navigation';
import { EmptyListCTA } from 'components/common';

const Schedule = () => {
  //const styles = useStyles2(getStyles);

  return (
    <EmptyListCTA
      title="You haven't created any schedule(s) yet."
      buttonTitle={'Create new schedule'}
      buttonIcon="calendar-alt"
      buttonLink={`${prefixRoute(ROUTES.SCHEDULERS)}/create`}
      proTip="Schedules are set with report groups and panels. They run on fixed schedules recursively sending excel reports through email."
    />
  );
};

// const getStyles = (theme: GrafanaTheme2) => ({
//   marginTop: css`
//     margin-top: ${theme.spacing(2)};
//   `,
// });

export { Schedule };
