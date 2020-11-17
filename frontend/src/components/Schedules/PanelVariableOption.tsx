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
  console.log(name, selectedOptions, selectableOptions);
  return (
    <div style={{ display: 'flex', flexDirection: 'row', marginTop: '5px' }}>
      <InlineFormLabel>{name}</InlineFormLabel>
      {!multiSelectable ? (
        <Select
          value={selectableOptions.filter((f: any) => !!selectedOptions?.find((s1: any) => s1 === f.value.value))}
          onChange={(selected: SelectableValue<SelectableVariable>) => onUpdate([selected])}
          options={selectableOptions}
        />
      ) : (
        <MultiSelect
          onChange={(selected: SelectableValue<SelectableVariable>) => onUpdate(selected)}
          value={selectableOptions.filter(
            (option: SelectableValue<SelectableVariable>) =>
              !!selectedOptions?.find((selected: string) => selected === option?.value?.value)
          )}
          filterOption={(option: SelectableValue<SelectableVariable>, searchQuery: string) =>
            !!option?.label?.toLowerCase().includes(searchQuery.toLowerCase())
          }
          closeMenuOnSelect={false}
          options={selectableOptions}
        />
      )}
    </div>
  );
};
