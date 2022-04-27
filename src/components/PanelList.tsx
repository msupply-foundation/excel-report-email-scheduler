import React, { useState } from 'react';
import { GrafanaTheme2 } from '@grafana/data';
import { EmptySearchResult, FieldSet, HorizontalGroup, Icon, Legend, Tag, Tooltip, useStyles2 } from '@grafana/ui';
import { css, cx } from '@emotion/css';
import { Panel } from 'types';
import intl from 'react-intl-universal';
import { useWindowSize } from '../hooks/useWindowResize';
import { PanelItem } from 'components';

const pageLimit = 20;

type PanelListProps = {
  panels: any[];
  panelListError: any;
  onPanelChecked: (event: React.FormEvent<HTMLInputElement>, panelID: number) => void;
  checkedPanels: number[];
};

const PanelList: React.FC<PanelListProps> = ({ panels, panelListError, onPanelChecked, checkedPanels }) => {
  const styles = useStyles2(getStyles);
  const { height } = useWindowSize();
  const [data, _] = useState<Panel[] | undefined>(panels);

  return (
    <>
      <div className="page-action-bar">
        <FieldSet label="Selected Panels">
          {checkedPanels.length > 0 ? (
            <HorizontalGroup wrap={true} style={{ marginBottom: '25px' }} align="flex-start" justify="flex-start">
              {checkedPanels.map((panelID) => {
                const panel = panels.find((panel: Panel) => panel.id === panelID);
                return <Tag key={panelID} icon="user" name={panel?.title} />;
              })}
            </HorizontalGroup>
          ) : (
            <EmptySearchResult>You have not selected any panels(s) yet</EmptySearchResult>
          )}
        </FieldSet>
      </div>
      <div style={{ marginTop: '25px' }}>
        <div style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
          <Tooltip placement="top" content={intl.get('available_panels_tooltip')} theme={'info'}>
            <Icon
              name="info-circle"
              size="sm"
              style={{ marginLeft: '10px', marginRight: '10px', marginBottom: '16px' }}
            />
          </Tooltip>
          <Legend>{intl.get('available_panels')}</Legend>
        </div>

        <ol className={styles.list} style={{ maxHeight: `${(height ?? 0) / 2}px`, overflow: 'scroll' }}>
          {data?.map((panel: any, key) => {
            return <PanelItem panel={panel} key={`panelItem${key}`} />;
          })}
        </ol>
      </div>
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  checkboxWrapper: css`
    label {
      line-height: 1.2;
    }
  `,
  list: cx('card-section', 'card-list-layout-grid', 'card-list-layout-list', 'card-list'),
});

export { PanelList };
