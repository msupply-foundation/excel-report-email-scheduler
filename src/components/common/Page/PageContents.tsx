// Libraries
import React, { FC } from 'react';
import { cx } from '@emotion/css';

interface Props {
  children: React.ReactNode;
  className?: string;
}

export const PageContents: FC<Props> = ({ children, className }) => {
  return (
    <div className={cx('page-container', 'page-body', className)} style={{ marginTop: '20px' }}>
      {children}
    </div>
  );
};
