import { SelectableValue } from '@grafana/data';
import { InlineFormLabel, MultiSelect, Select } from '@grafana/ui';
import { refreshPanelOptions } from 'api';
import { useDatasourceID } from 'hooks';
import React from 'react';
import { useQuery } from 'react-query';
import { SelectableVariable, Variable } from 'types';

type Props = {
  onUpdate: (selected: SelectableValue) => void;
  name: string;
  multiSelectable: boolean;

  selectableOptions: Array<SelectableValue<SelectableVariable>>;
  variable: Variable;
};

export const PanelVariableOptions: React.FC<Props> = ({
  onUpdate,
  name,
  multiSelectable,
  selectableOptions,
  variable,
}) => {
  const { refresh } = variable;
  const datasourceID = useDatasourceID();

  // When a query variable is set to refresh, it does not by default have `options` pre-populated.
  // So, when refresh is true, query for the data and map it to the matching array.
  const { data } = useQuery(name, () => refreshPanelOptions(variable, datasourceID), {
    enabled: !!refresh,
  });

  const options = selectableOptions?.length > 0 ? selectableOptions : data;

  return (
    <div style={{ display: 'flex', flexDirection: 'row', marginTop: '5px' }}>
      <InlineFormLabel>{name}</InlineFormLabel>
      {multiSelectable ? (
        <Select
          //value={options?.filter((f: any) => !!selectedOptions?.find((s1: any) => s1 === f.value.value))}
          onChange={(selected: SelectableValue<SelectableVariable>) => {
            if (selected.value?.value === '$__all') {
              return onUpdate(selectableOptions);
            }

            onUpdate([selected]);
          }}
          options={options}
        />
      ) : (
        <MultiSelect
          onChange={(selected: SelectableValue<SelectableVariable>) => {
            const isAll = selected.some(({ value }: SelectableValue<SelectableVariable>) => value?.value === '$__all');
            if (isAll) {
              return onUpdate(selectableOptions);
            }

            onUpdate(selected);
          }}
          closeMenuOnSelect={false}
          options={options}
        />
      )}
    </div>
  );
};
