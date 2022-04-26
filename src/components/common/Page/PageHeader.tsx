import React, { FC } from 'react';
import { css } from '@emotion/css';

import { GrafanaTheme2 } from '@grafana/data';
import { IconName, LinkButton, useStyles2 } from '@grafana/ui';
import { HeaderPropsType } from './Page';

const PageHeader: FC<HeaderPropsType> = (props) => {
  const styles = useStyles2(getStyles);

  return (
    <div className={styles.headerCanvas}>
      <div className="page-container">
        <div className="page-header">{renderHeaderTitle(props)}</div>
        {props.backButton && (
          <LinkButton icon={props.backButton.icon as IconName} href={props.backButton.href}>
            Back
          </LinkButton>
        )}
      </div>
    </div>
  );
};

function renderHeaderTitle({ title, subTitle }: HeaderPropsType) {
  return (
    <div className="page-header__inner">
      <div className="page-header__info-block">
        <h1 className="page-header__title">{title}</h1>
        {subTitle && <div className="page-header__sub-title">{subTitle}</div>}
      </div>
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2) => ({
  headerCanvas: css`
    background: ${theme.colors.background.canvas};
  `,
});

export { PageHeader };
