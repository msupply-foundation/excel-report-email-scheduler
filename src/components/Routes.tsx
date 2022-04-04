import { ROUTES } from '../constants';
import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { prefixRoute, useNavigation } from 'routes/router';
import { PageOne } from './Pages/PageOne';

export const AppRouter = () => {
  useNavigation();

  return (
    <Routes>
      <Route path={prefixRoute(ROUTES.One)} element={PageOne} />
    </Routes>
  );
};
