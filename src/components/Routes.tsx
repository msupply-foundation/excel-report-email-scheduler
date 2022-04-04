import { ROUTES } from '../constants';
import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { prefixRoute, useNavigation } from 'routes/router';
import { ReportGroup, Scheduler } from '../Pages';

export const AppRouter = () => {
  useNavigation();

  return (
    <Routes>
      <Route path={prefixRoute(ROUTES.REPORT_GROUP)} element={ReportGroup} />
      <Route path={prefixRoute(ROUTES.SCHEDULERS)} element={Scheduler} />
    </Routes>
  );
};
