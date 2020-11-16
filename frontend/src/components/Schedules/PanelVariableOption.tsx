import { SelectableValue } from '@grafana/data';
import { InlineFormLabel, MultiSelect, Select } from '@grafana/ui';
import { SelectableVariable } from 'common/types';
import React, { FC } from 'react';

type Props = {
  onUpdate: (selected: SelectableValue) => void;
  name: string;
  multiSelectable: boolean;
  selectedOptions: string[];
  selectableOptions: Array<SelectableValue<SelectableVariable>>;
};

export const PanelVariableOptions: FC<Props> = ({
  onUpdate,
  name,
  multiSelectable,
  selectedOptions,
  selectableOptions,
}) => {
  return (
    <div style={{ display: 'flex', flexDirection: 'row', marginTop: '5px' }}>
      <InlineFormLabel>{name}</InlineFormLabel>
      {!multiSelectable ? (
        <Select
          value={selectableOptions.filter((f: any) => !!selectedOptions?.find((s1: any) => s1 === f.value.value))}
          onChange={(selectedLookback: SelectableValue) => onUpdate([selectedLookback])}
          options={selectableOptions}
        />
      ) : (
        <MultiSelect
          onChange={(selectedLookback: SelectableValue) => onUpdate(selectedLookback)}
          value={selectableOptions.filter((f: any) => !!selectedOptions?.find((s1: any) => s1 === f.value.value))}
          filterOption={(option: SelectableValue, searchQuery: string) =>
            !!option?.text?.startsWith(searchQuery.toLowerCase())
          }
          closeMenuOnSelect={false}
          options={selectableOptions}
        />
      )}
    </div>
  );
};
