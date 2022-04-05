import { ROUTES } from '../constants';
import { Redirect, Route, Switch } from 'react-router-dom';
import { prefixRoute, useNavigation } from '../utils/navigation';
import { ReportGroup, Scheduler } from '../Pages';
import React from 'react';

const AppRoutes = () => {
  useNavigation();

  return (
    <Switch>
      <Route exact path={prefixRoute(ROUTES.REPORT_GROUP)} component={ReportGroup} />
      <Route exact path={prefixRoute(ROUTES.SCHEDULERS)} component={Scheduler} />

      {/* Default page */}
      <Route exact path="*">
        <Redirect to={prefixRoute(ROUTES.REPORT_GROUP)} />
      </Route>
    </Switch>
  );
};

export { AppRoutes };
