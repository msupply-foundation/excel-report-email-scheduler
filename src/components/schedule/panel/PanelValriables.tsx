import React from 'react';
import { SelectableValue } from '@grafana/data';
import { Tooltip, Icon, InlineFormLabel, Select } from '@grafana/ui';
import { getLookbacks } from '../../../constants';

import intl from 'react-intl-universal';
import { ContentVariables, Panel, PanelDetails, SelectableVariable, Variable, VariableOption } from 'types';
import { PanelVariableOptions } from './PanelVariableOptions';
import { panelUsesMacro, parseOrDefault } from 'utils';
import { PanelVariableTextInput } from 'components';

type Props = {
  panel: Panel;
  panelDetail: PanelDetails;
  onUpdateLookback: (selectedValue: SelectableValue) => void;
  onUpdateVariable: (variableName: string) => (selectedValue: SelectableValue) => void;
};

export const PanelVariables: React.FC<Props> = ({ panel, onUpdateVariable, panelDetail, onUpdateLookback }) => {
  const lookbacks = getLookbacks();

  const vars = parseOrDefault<ContentVariables>(panelDetail?.variables, {});

  const usesMacro = panelUsesMacro(panel.rawSql);
  const usesVariables = panel.variables.length > 0;

  if (!(usesVariables || usesMacro)) {
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
            value={!!panelDetail && panelDetail.lookback}
            onChange={(selected: SelectableValue<String>) => onUpdateLookback(selected)}
          />
        </div>
      )}

      {panel.variables.map((variable: Variable) => {
        const { name, options: variableOptions, multi, label } = variable;

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

        if (variable.type === 'textbox') {
          // Pre-fill with either the report content value that has been saved in the msupply sqlite,
          // or what is currently being used in the dashboard, as a default.
          const value = selected?.[0] ?? options[0]?.value?.value;
          return <PanelVariableTextInput onUpdate={onUpdateVariable(name)} name={label ?? name} value={value} />;
        }

        return (
          <PanelVariableOptions
            onUpdate={onUpdateVariable(name)}
            key={name}
            multiSelectable={multi}
            name={label ?? name}
            variable={variable}
            selectedOptions={selected}
            selectableOptions={options}
          />
        );
      })}
    </div>
  );
};
