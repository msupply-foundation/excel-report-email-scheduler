import React, { FC, useEffect } from 'react';
import intl from 'react-intl-universal';
import { ReactQueryDevtools } from 'react-query-devtools';
import { AppRootProps } from '@grafana/data';
import { locales } from '../common/translations';
import { useTabs } from 'hooks';
import { ReportGroupTab } from './ReportGroupTab';
import { ReportSchedulesTab } from './ReportSchedulesTab';

interface Props extends AppRootProps {}

export const TAB_ID_REPORT_GROUP = 'report-group';
export const TAB_ID_REPORT_SCHEDULE = 'report-schedule';

export const RootPage: FC<Props> = ({ path, onNavChanged, query, meta }) => {
  const pathWithoutLeadingSlash = path.replace(/^\//, '');

  useEffect(() => {
    intl.init({ currentLocale: 'en', locales });
  }, []);

  useTabs(query.tab, meta.info.logos.large, pathWithoutLeadingSlash, onNavChanged);

  const getTabContent = () => {
    switch (query.tab) {
      default:
      case TAB_ID_REPORT_GROUP: {
        return <ReportGroupTab meta={meta} path={path} onNavChanged={onNavChanged} query={query} />;
      }
      case TAB_ID_REPORT_SCHEDULE: {
        return <ReportSchedulesTab meta={meta} path={path} onNavChanged={onNavChanged} query={query} />;
      }
    }
  };

  return (
    <div>
      {getTabContent()}
      <ReactQueryDevtools />
    </div>
  );
};
