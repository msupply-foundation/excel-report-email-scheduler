// Libraries
import React, { FC, HTMLAttributes } from 'react';

import { CustomScrollbar, useStyles2 } from '@grafana/ui';
import { GrafanaTheme2 } from '@grafana/data';
import { css, cx } from '@emotion/css';
import { PageHeader } from './PageHeader';
import { PageContents } from './PageContents';

interface Props extends HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  headerProps?: HeaderPropsType;
}

export interface HeaderPropsType {
  title: string;
  subTitle: string;
  backButton?: {
    icon: string;
    href: string;
  };
}

export interface PageType extends FC<Props> {
  Contents: typeof PageContents;
}

export const Page: PageType = ({ headerProps, children, className, ...otherProps }) => {
  const styles = useStyles2(getStyles);

  return (
    <div {...otherProps} className={cx(styles.wrapper, className)}>
      <CustomScrollbar autoHeightMin={'100%'}>
        <div className="page-scrollbar-content">
          {headerProps && <PageHeader {...headerProps} />}
          {children}
        </div>
      </CustomScrollbar>
    </div>
  );
};

Page.Contents = PageContents;

const getStyles = (theme: GrafanaTheme2) => ({
  wrapper: css`
    width: 100%;
    flex-grow: 1;
    min-height: 0;
  `,
});
