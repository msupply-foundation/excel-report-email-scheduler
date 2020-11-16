import { useEffect } from 'react';
import { TAB_ID_REPORT_GROUP, TAB_ID_REPORT_SCHEDULE } from 'components/RootPage';
import { NavModel } from '@grafana/data';
import intl from 'react-intl-universal';

export const useTabs = (activeTab: string, logo: string, path: string, onNavChanged: (nav: NavModel) => void) => {
  useEffect(() => {
    const tabs = [
      {
        text: `    ${intl.get('report_groups')}`,
        icon: 'fa fa-fw fa-file-text-o',
        url: path + '?tab=' + TAB_ID_REPORT_GROUP,
        id: TAB_ID_REPORT_GROUP,
      },
      {
        text: `    ${intl.get('report_schedules')}`,
        icon: 'fa fa-users',
        url: path + '?tab=' + TAB_ID_REPORT_SCHEDULE,
        id: TAB_ID_REPORT_SCHEDULE,
      },
    ];

    const node = {
      text: intl.get('msupply'),
      img: logo,
      url: path,
      children: tabs,
    };

    onNavChanged({
      node: node,
      main: node,
    });
  }, [activeTab, path, logo, onNavChanged]);
};
