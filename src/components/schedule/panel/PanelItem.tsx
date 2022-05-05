import React, { useContext } from 'react';
import { ContentVariables, Panel, PanelDetails, SelectableVariable } from 'types';
import { GrafanaTheme2, SelectableValue } from '@grafana/data';
import { Checkbox, useStyles2 } from '@grafana/ui';
import { css } from '@emotion/css';
import { PanelContext } from 'context';
import { PanelVariables } from 'components';
import { parseOrDefault } from 'utils';
//import { PanelVariablesContext } from 'context';

type Props = {
  panel: Panel;
  onPanelChecked: (panel: Panel) => void;
  checkedPanels: number[];
  panelDetail: PanelDetails;
};

export const PanelItem: React.FC<Props> = ({ panel, onPanelChecked, panelDetail, checkedPanels }) => {
  const styles = useStyles2(getStyles);
  const { title, description, error } = panel;

  const { setPanelDetails } = useContext(PanelContext);

  const onUpdateLookback = (content: PanelDetails) => (selectableValue: SelectableValue) => {
    setPanelDetails((prevPanels: any) => {
      const myIndex = prevPanels.findIndex(
        (el: any) => el.panelID === content.panelID && el.dashboardID === content.dashboardID
      );

      return [
        ...prevPanels.slice(0, myIndex),
        { ...prevPanels[myIndex], lookback: selectableValue.value },
        ...prevPanels.slice(myIndex + 1),
      ];
    });
  };

  const onUpdateVariable =
    (content: PanelDetails) => (variableName: string) => (selectableValue: SelectableValue<SelectableVariable[]>) => {
      const newVariable = selectableValue.map(({ value }: SelectableValue) => value.value);
      const newVariables = parseOrDefault<ContentVariables>(content.variables, {});

      newVariables[variableName] = newVariable;

      setPanelDetails((prevPanels: any) => {
        const myIndex = prevPanels.findIndex(
          (el: any) => el.panelID === panel.id && el.dashboardID === panel.dashboardID
        );

        return [
          ...prevPanels.slice(0, myIndex),
          { ...prevPanels[myIndex], variables: JSON.stringify(newVariables) },
          ...prevPanels.slice(myIndex + 1),
        ];
      });
    };

  return (
    <li className="card-item-wrapper" style={{ cursor: !error ? 'pointer' : '' }}>
      <div className="card-item">
        <div
          className="card-item-body"
          onClick={(event) => {
            event.preventDefault();
            if (error) {
              return;
            }
            return onPanelChecked(panel);
          }}
        >
          <div className={styles.marginForCheckbox}>
            {!error ? (
              <Checkbox
                value={!!checkedPanels && !!checkedPanels.some((checkedPanelID: number) => checkedPanelID === panel.id)}
              />
            ) : null}
          </div>
          <div className="card-item-details">
            <div className="card-item-name">{title}</div>
            <div className="card-item-type">{description}</div>
          </div>
        </div>

        <PanelVariables
          panel={panel}
          checkedPanels={checkedPanels}
          panelDetail={panelDetail}
          onUpdateVariable={onUpdateVariable(panelDetail)}
          onUpdateLookback={onUpdateLookback(panelDetail)}
        />
      </div>
    </li>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  marginForCheckbox: css`
    margin-right: 10px;
  `,
});
