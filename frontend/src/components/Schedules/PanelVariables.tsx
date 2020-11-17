import React, { FC } from 'react';
import { Icon, InlineFormLabel, Select, Tooltip } from '@grafana/ui';
import intl from 'react-intl-universal';
import { SelectableValue } from '@grafana/data';
import { ReportContentKey } from 'common/enums';
import { ContentVariables, Panel, SelectableVariable, Store, Variable, VariableOption } from 'common/types';

import { getLookbacks, parseOrDefault } from 'common';
import { PanelVariableOptions } from './PanelVariableOption';
import { panelUsesMacro } from 'common/utils/checkers';

type Props = {
  storeIDs: string;
  lookback: number;
  variables: string;
  panel: Panel;
  onUpdateVariable: (variableName: string) => (selectedValue: SelectableValue) => void;
  onUpdateContent: (key: ReportContentKey, selectedValue: SelectableValue<String | Number | Store>) => void;
};

export const PanelVariables: FC<Props> = ({ onUpdateVariable, panel, onUpdateContent, lookback, variables }) => {
  const lookbacks = getLookbacks();
  const vars = parseOrDefault<ContentVariables>(variables, {});

  const usesMacro = panelUsesMacro(panel.rawSql);
  const usesVariables = usesMacro && panel.variables.length > 0;
  if (!(usesVariables && usesMacro)) {
    return null;
  }

  return (
    <div style={{ border: '1px solid grey', padding: '20px' }}>
      <div className="card-item-type">{intl.get('variables')}</div>
      <Tooltip placement="top" content={intl.get('variables_tooltip')} theme={'info'}>
        <Icon name="info-circle" size="sm" style={{ marginLeft: '10px' }} />
      </Tooltip>

      {usesMacro && (
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
      )}

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
    </div>
  );
};
