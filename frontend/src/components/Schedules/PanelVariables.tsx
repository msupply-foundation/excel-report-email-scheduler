import React, { FC } from 'react';
import { Icon, InlineFormLabel, MultiSelect, Select, Tooltip } from '@grafana/ui';
import intl from 'react-intl-universal';
import { SelectableValue } from '@grafana/data';
import { ReportContentKey } from 'common/enums';
import { ContentVariables, Panel, SelectableVariable, Store, Variable, VariableOption } from 'common/types';
import { useStores } from 'hooks/useStores';
import { getLookbacks, parseOrDefault } from 'common';
import { PanelVariableOptions } from './PanelVariableOption';
import { panelUsesVariable } from 'common/utils/checkers';

type Props = {
  storeIDs: string;
  lookback: number;
  variables: string;
  panel: Panel;
  onUpdateVariable: (variableName: string) => (selectedValue: SelectableValue) => void;
  onUpdateContent: (key: ReportContentKey, selectedValue: SelectableValue<String | Number | Store>) => void;
};

export const PanelVariables: FC<Props> = ({
  onUpdateVariable,
  panel,
  onUpdateContent,
  storeIDs,
  lookback,
  variables,
}) => {
  const lookbacks = getLookbacks();
  const stores = useStores();
  const vars = parseOrDefault<ContentVariables>(variables, {});

  const usesStores = panelUsesVariable(panel.rawSql, 'store');

  return (
    <>
      <div className="card-item-type">{intl.get('variables')}</div>
      <Tooltip placement="top" content={intl.get('variables_tooltip')} theme={'info'}>
        <Icon name="info-circle" size="sm" style={{ marginLeft: '10px' }} />
      </Tooltip>

      {usesStores && (
        <div style={{ display: 'flex', flexDirection: 'row' }}>
          <InlineFormLabel tooltip={intl.get('selected_stores_description')}>
            {intl.get('selected_stores')}
          </InlineFormLabel>
          <MultiSelect
            placeholder={intl.get('choose_stores')}
            closeMenuOnSelect={false}
            filterOption={(option: SelectableValue<Store>, searchQuery: string) =>
              !!option.label?.toLowerCase().startsWith(searchQuery.toLowerCase())
            }
            value={stores.filter(({ id }) => storeIDs.includes(id))}
            onChange={(selectedStores: SelectableValue<Store>) =>
              onUpdateContent(ReportContentKey.STORE_ID, selectedStores)
            }
            options={stores.map((store: any) => ({ label: store.name, value: store }))}
          />
        </div>
      )}

      <div style={{ display: 'flex', flexDirection: 'row' }}>
        <InlineFormLabel tooltip={intl.get('lookback_period_description')}>
          {intl.get('lookback_period')}
        </InlineFormLabel>
        <Select
          options={lookbacks}
          onChange={(selected: SelectableValue<Number>) => onUpdateContent(ReportContentKey.LOOKBACK, selected)}
          value={lookbacks.filter((selected: SelectableValue<Number>) => selected.value === lookback)}
        />
      </div>

      {panel.variables.map((variable: Variable) => {
        const { name, options: variableOptions, multi } = variable;
        // For a panels variables, find the ones which are selected from the
        // ReportContent.variables field, which is a stringified object consisting
        // of { [variable.name]: [Array of chosen options as strings] }
        // For example, a variable ${VEN} could have the `ReportContent.variables`
        // field { VEN: ['V', 'E'] }.
        const selected = vars[variable.name];
        // As well as mapping the current available options into a SelectedValue to
        const options: Array<SelectableValue<SelectableVariable>> = variableOptions.map((option: VariableOption) => {
          return { label: option.text, value: { name: variable.name, value: option.value } };
        });

        return (
          <PanelVariableOptions
            onUpdate={onUpdateVariable(name)}
            multiSelectable={multi}
            name={name}
            selectedOptions={selected}
            selectableOptions={options}
          />
        );
      })}
    </>
  );
};
