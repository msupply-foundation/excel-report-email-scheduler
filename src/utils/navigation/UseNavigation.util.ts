import { NavModel } from '@grafana/data';
import { NAVIGATION } from '../../constants';
import { usePluginProps } from 'context';
import { useLocation } from 'react-router-dom';
import { getNavModel } from '.';
import { useEffect } from 'react';

// Displays a top navigation tab-bar if needed
function useNavigation() {
  const pluginProps = usePluginProps();
  const location = useLocation();

  useEffect(() => {
    const excludeURIs: string[] = ['create', 'edit'];

    if (!pluginProps) {
      console.error('Root plugin props are not available in the context.');
      return;
    }

    const activeId = excludeURIs.some((uri) => location.pathname.includes(uri))
      ? ''
      : Object.keys(NAVIGATION).find((routeId) => location.pathname.includes(routeId)) || '';
    const activeNavItem = NAVIGATION[activeId];

    const { onNavChanged, meta, basename } = pluginProps;

    // Disable tab navigation
    // (the route is not registered as a navigation item)
    if (!activeNavItem) {
      onNavChanged(undefined as unknown as NavModel);
    }

    // Show tabbed navigation with the active tab
    else {
      onNavChanged(
        getNavModel({
          activeId,
          basePath: basename,
          logoUrl: meta.info.logos.large,
        })
      );
    }
  }, [location.pathname, pluginProps]);
}

export { useNavigation };
