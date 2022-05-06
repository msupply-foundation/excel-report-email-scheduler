import React, { useContext } from 'react';
import { GrafanaTheme2 } from '@grafana/data';
import { Card, FieldSet, HorizontalGroup, Icon, Legend, Tag, Tooltip, useStyles2 } from '@grafana/ui';
import { css, cx } from '@emotion/css';
import { Panel, PanelDetails } from 'types';
import intl from 'react-intl-universal';
import { Loading, PanelItem } from 'components';
import { PanelContext } from 'context';

//const pageLimit = 20;

type PanelListProps = {
  panelListError: any;
  onPanelChecked: (panel: Panel) => void;
  checkedPanels: number[];
};

const PanelList: React.FC<PanelListProps> = ({ panelListError, onPanelChecked, checkedPanels }) => {
  const styles = useStyles2(getStyles);

  const { panels, panelDetails } = useContext(PanelContext);

  if (!panels) {
    return <Loading />;
  }

  return (
    <>
      <div className="page-action-bar">
        <Card disabled={true}>
          <Card.Heading>Selected panel(s)</Card.Heading>
          <Card.Figure>
            <Icon size="xxxl" name="bell" />
          </Card.Figure>
          <Card.Description>
            <FieldSet>
              {!!checkedPanels && checkedPanels.length > 0 ? (
                <HorizontalGroup wrap={true} style={{ marginBottom: '25px' }} align="flex-start" justify="flex-start">
                  {checkedPanels.map((checkedPanel) => {
                    const panel = panels.find((panel: Panel) => panel.id === checkedPanel);
                    if (!panel) {
                      return false;
                    }

                    return <Tag key={checkedPanel} icon="user" name={panel.title} />;
                  })}
                </HorizontalGroup>
              ) : (
                <div>You have not selected any panels(s) yet. Please select one or more panels from the list below</div>
              )}
            </FieldSet>
          </Card.Description>
        </Card>
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
        <ol className={styles.list}>
          {panels &&
            panelDetails &&
            panels?.map((panel: Panel, key: any) => (
              <PanelItem
                panel={panel}
                key={`panelItem${key}`}
                onPanelChecked={onPanelChecked}
                checkedPanels={checkedPanels}
                panelDetail={
                  panelDetails.find(
                    (detail: PanelDetails) => detail.panelID === panel.id && detail.dashboardID === panel.dashboardID
                  )!
                }
              />
            ))}
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
