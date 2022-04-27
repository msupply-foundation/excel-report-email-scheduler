import React from 'react';
import { Panel } from 'types';
import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';
import { PanelVariables } from 'components';

type Props = {
  panel: Panel;
  //onToggle: (panel: Panel) => Promise<void>;
};

export const PanelItem: React.FC<Props> = ({ panel }) => {
  const styles = useStyles2(getStyles);
  const { title, description, error } = panel;

  return (
    <li className="card-item-wrapper" style={{ cursor: !error ? 'pointer' : '' }}>
      <div className={'card-item'}>
        <div className="card-item-body">
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

const getStyles = (theme: GrafanaTheme2) => ({});
