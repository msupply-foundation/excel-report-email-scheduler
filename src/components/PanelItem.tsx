import React from 'react';
import { Panel } from 'types';
import { GrafanaTheme2 } from '@grafana/data';
import { Checkbox, useStyles2 } from '@grafana/ui';
import { PanelVariables } from 'components';
import { css } from '@emotion/css';

type Props = {
  panel: Panel;
  onToggle: (panel: Panel) => void;
  checkedPanels: number[];
};

export const PanelItem: React.FC<Props> = ({ panel, onToggle, checkedPanels }) => {
  const styles = useStyles2(getStyles);
  const { title, description, error } = panel;

  return (
    <li className="card-item-wrapper" style={{ cursor: !error ? 'pointer' : '' }}>
      <div className="card-item">
        <div
          className="card-item-body"
          onClick={() => {
            if (error) {
              return;
            }
            return onToggle(panel);
          }}
        >
          <div className={styles.marginForCheckbox}>
            {!error ? <Checkbox value={checkedPanels.includes(panel.id)} /> : null}
          </div>
          <div className="card-item-details">
            <div className="card-item-name">{title}</div>
            <div className="card-item-type">{description}</div>
          </div>
        </div>

        <PanelVariables panel={panel} />
      </div>
    </li>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  marginForCheckbox: css`
    margin-right: 10px;
  `,
});
