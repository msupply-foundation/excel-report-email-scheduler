import React, { FC } from 'react';

import { InlineFormLabel, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

type Props = {
  tooltip: string;
  value: any;
  options: SelectableValue[];
  onChange: (selected: SelectableValue) => void;
  label: string;
};

export const FieldSelect: FC<Props> = ({ tooltip, value, options, onChange, label }) => {
  return (
    <div className="gf-form">
      <InlineFormLabel className="width-14" tooltip={tooltip}>
        {label}
      </InlineFormLabel>
      <div style={{ display: 'flex', flex: 1, flexDirection: 'column' }}>
        <Select value={value} options={options} onChange={onChange} />
      </div>
    </div>
  );
};
