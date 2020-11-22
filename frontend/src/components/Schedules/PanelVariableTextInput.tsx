import { css } from 'emotion';
import { SelectableValue } from '@grafana/data';
import { InlineFormLabel, Input } from '@grafana/ui';

import React, { FC } from 'react';

const flexContainer = css`
  display: flex;
  flex: 1;
  flex-basis: 50%;
`;

type Props = {
  onUpdate: (selected: SelectableValue) => void;
  name: string;
  value: string;
};

export const PanelVariableTextInput: FC<Props> = ({ onUpdate, name, value }) => {
  return (
    <div style={{ display: 'flex', flexDirection: 'row', marginTop: '5px', flexWrap: 'wrap' }}>
      <InlineFormLabel>{name}</InlineFormLabel>
      <div className={flexContainer}>
        <Input
          type="text"
          defaultValue={value}
          placeholder=""
          onChange={e => {
            const { value } = e.target as HTMLInputElement;
            onUpdate([{ value: { name, value } }]);
          }}
          css=""
        />
      </div>
    </div>
  );
};
