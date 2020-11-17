import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { Button } from '@grafana/ui';
import { css } from 'emotion';

type DisabledProps = {
  toggle: Function;
};

export const DisabledState: FC<DisabledProps> = ({ toggle }) => {
  const centered = css`
    align-items: center;
    justify-content: center;
    flex: 1;
    display: flex;
    height: 100%;
  `;

  return (
    <div className={centered}>
      <Button variant="primary" onClick={() => toggle(true)}>
        {intl.get('enable')}
      </Button>
    </div>
  );
};
