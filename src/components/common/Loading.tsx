import React, { HTMLAttributes } from 'react';
import { css } from '@emotion/css';
import { LoadingPlaceholder, useStyles2 } from '@grafana/ui';

export interface LoadingPlaceholderProps extends HTMLAttributes<HTMLDivElement> {
  text?: React.ReactNode;
}

const Loading: React.FC<LoadingPlaceholderProps> = ({ text = 'loading...', ...rest }: LoadingPlaceholderProps) => {
  const style = useStyles2(getStyles);

  return (
    <div className={style.loadingWrapper}>
      <LoadingPlaceholder text={text} {...rest} />
    </div>
  );
};

const getStyles = () => ({
  loadingWrapper: css`
    display: flex;
    height: 50vh;
    align-items: center;
    justify-content: center;
  `,
});

export { Loading };
