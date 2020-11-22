import { SelectableValue } from '@grafana/data';
import { css } from 'emotion';
import { InlineFormLabel, MultiSelect, Select } from '@grafana/ui';
import { refreshPanelOptions } from 'api';
import { SelectableVariable, Variable } from 'common/types';
import { useDatasourceID } from 'hooks';
import React, { FC } from 'react';
import { useQuery } from 'react-query';

const flexContainer = css`
  display: flex;
  flex: 1;
  flex-basis: 50%;
`;

type Props = {
  onUpdate: (selected: SelectableValue) => void;
  name: string;
  multiSelectable: boolean;
  selectedOptions: string[];
  selectableOptions: Array<SelectableValue<SelectableVariable>>;
  variable: Variable;
};

export const PanelVariableOptions: FC<Props> = ({
  onUpdate,
  name,
  variable,
  multiSelectable,
  selectedOptions,
  selectableOptions,
}) => {
  const { refresh } = variable;
  const datasourceID = useDatasourceID();
  console.log(selectedOptions);
  // When a query variable is set to refresh, it does not by default have `options` pre-populated.
  // So, when refresh is true, query for the data and map it to the matching array.
  const { data } = useQuery(name, () => refreshPanelOptions(variable, datasourceID), {
    enabled: !!refresh,
  });

  const options = selectableOptions?.length > 0 ? selectableOptions : data;

  return (
    <div style={{ display: 'flex', flexDirection: 'row', marginTop: '5px', flexWrap: 'wrap' }}>
      <InlineFormLabel>{name}</InlineFormLabel>
      <div className={flexContainer}>
        {!multiSelectable ? (
          <Select
            value={options?.filter((f: any) => !!selectedOptions?.find((s1: any) => s1 === f.value.value))}
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
              const isAll = selected.some(
                ({ value }: SelectableValue<SelectableVariable>) => value?.value === '$__all'
              );
              if (isAll) {
                return onUpdate(selectableOptions);
              }

              onUpdate(selected);
            }}
            value={options?.filter(
              (option: SelectableValue<SelectableVariable>) =>
                !!selectedOptions?.find((selected: string) => selected === option?.value?.value)
            )}
            filterOption={(option: SelectableValue<SelectableVariable>, searchQuery: string) =>
              !!option?.label?.toLowerCase().includes(searchQuery.toLowerCase())
            }
            closeMenuOnSelect={false}
            options={options}
          />
        )}
      </div>
    </div>
  );
};
