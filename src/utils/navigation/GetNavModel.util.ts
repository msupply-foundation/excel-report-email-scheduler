import { NAVIGATION_TITLE, NAVIGATION_SUBTITLE, NAVIGATION } from '../../constants';

function getNavModel({ activeId, basePath, logoUrl }: { activeId: string; basePath: string; logoUrl: string }) {
  const main = {
    text: NAVIGATION_TITLE,
    subTitle: NAVIGATION_SUBTITLE,
    url: basePath,
    img: logoUrl,
    children: Object.values(NAVIGATION).map((navItem) => ({
      ...navItem,
      active: navItem.id === activeId,
    })),
  };

  return {
    main,
    node: main,
  };
}

export { getNavModel };
