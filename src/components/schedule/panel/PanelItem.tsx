import React, { useContext } from 'react';
import { Panel, PanelDetails, PanelListSelectedType } from 'types';
import { GrafanaTheme2 } from '@grafana/data';
import { Checkbox, useStyles2 } from '@grafana/ui';
import { css } from '@emotion/css';
import { PanelContext } from 'context';
import { PanelVariables } from 'components';

type Props = {
  panel: Panel;
  onPanelChecked: (panel: Panel) => void;
  checkedPanels: PanelListSelectedType[];
  panelDetail: PanelDetails;
};

export const PanelItem: React.FC<Props> = ({ panel, onPanelChecked, panelDetail, checkedPanels }) => {
  const styles = useStyles2(getStyles);
  const { title, description, error } = panel;

  const { onUpdateLookback, onUpdateVariable } = useContext(PanelContext);

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
                value={
                  !!checkedPanels &&
                  !!checkedPanels.some((checkedPanel: PanelListSelectedType) => {
                    return checkedPanel.panelID === panel.id && checkedPanel.dashboardID === panel.dashboardID;
                  })
                }
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
          panelDetail={panelDetail}
          onUpdateVariable={onUpdateVariable(panelDetail, panel)}
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
