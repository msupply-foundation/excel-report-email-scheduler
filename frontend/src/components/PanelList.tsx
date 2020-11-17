import React, { FC } from 'react';
import classNames from 'classnames';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Checkbox, Field, MultiSelect, Select } from '@grafana/ui';
import {
  createReportContent,
  deleteReportContent,
  getPanels,
  getReportContent,
  getStores,
  updateReportContent,
} from '../api';
import { SelectableValue } from '@grafana/data';
import { queryCache, useMutation, useQuery } from 'react-query';
import { Schedule, ReportContent } from 'common/types';
import { ReportContentKey } from 'common/enums';

const getLookbacks = () => [
  { label: intl.get('1day'), value: 24 * 60 * 1000 * 60 },
  { label: intl.get('2days'), value: 2 * 24 * 60 * 1000 * 60 },
  { label: intl.get('3days'), value: 3 * 24 * 60 * 1000 * 60 },
  { label: intl.get('1week'), value: 7 * 24 * 60 * 1000 * 60 },
  { label: intl.get('2weeks'), value: 14 * 24 * 60 * 1000 * 60 },
  { label: intl.get('4weeks'), value: 28 * 24 * 60 * 1000 * 60 },
  { label: intl.get('3months'), value: 28 * 3 * 24 * 60 * 1000 * 60 },
  { label: intl.get('6months'), value: 28 * 6 * 24 * 60 * 1000 * 60 },
  { label: intl.get('1year'), value: 356 * 24 * 60 * 1000 * 60 },
];

const listStyle = classNames({
  'card-section': true,
  'card-list-layout-grid': true,
  'card-list-layout-list': true,
  'card-list': true,
});

const marginForCheckbox = css`
  margin-right: 10px;
`;

type Props = {
  schedule: Schedule;
  datasourceID: number;
};

// TODO: Create PanelItem component and extract logic there.
export const PanelList: FC<Props> = ({ schedule, datasourceID }) => {
  const { id: scheduleID } = schedule;

  const { data: panels } = useQuery(['panels'], getPanels);
  const { data: stores } = useQuery(['stores'], () => getStores(datasourceID));
  const { data: content } = useQuery<ReportContent[]>(['reportContent', scheduleID], getReportContent);

  const refetchContent = () => queryCache.refetchQueries(['reportContent', scheduleID]);
  const [createContent] = useMutation(createReportContent, { onSuccess: refetchContent });
  const [deleteContent] = useMutation(deleteReportContent, { onSuccess: refetchContent });
  const [updateContent] = useMutation(updateReportContent, { onSuccess: refetchContent });

  const onTogglePanel = async (content: ReportContent | null, panel: any) => {
    if (content) {
      await deleteContent(content);
    } else {
      await createContent({ scheduleID, panelID: panel?.id, dashboardID: panel?.dashboardID });
    }
  };

  const onUpdateContent = (
    content: ReportContent,
    key: ReportContentKey,
    selectableValue: SelectableValue<string | number>
  ) => {
    let newValue = selectableValue.value;

    if (key === ReportContentKey.STORE_ID && Array.isArray(selectableValue)) {
      newValue = selectableValue.map((store: any) => store?.value?.id).join(', ');
    }

    const newState = { ...content, [key]: newValue };
    updateContent(newState);
  };

  console.log(panels);

  return (
    <ol className={listStyle}>
      {panels?.map((panel: any) => {
        const { title, description } = panel;
        const matchingContent: ReportContent | null =
          content?.find(
            (reportContent: ReportContent) =>
              reportContent.panelID === panel.id && reportContent.dashboardID === panel.dashboardID
          ) ?? null;

        return (
          <li className="card-item-wrapper" style={{ cursor: 'pointer' }}>
            <div className={'card-item'}>
              <div className="card-item-body" onClick={() => onTogglePanel(matchingContent, panel)}>
                <div className={marginForCheckbox}>
                  <Checkbox value={!!matchingContent} css="" />
                </div>

                <div className="card-item-details">
                  <div className="card-item-name">{title}</div>
                  <div className="card-item-type">{description}</div>
                </div>
              </div>

              {matchingContent && (
                <>
                  <Field label={intl.get('selected_stores')} description={intl.get('selected_stores_description')}>
                    <MultiSelect
                      disabled={!matchingContent}
                      placeholder={matchingContent ? intl.get('choose_stores') : ''}
                      closeMenuOnSelect={false}
                      filterOption={(option: SelectableValue, searchQuery: string) =>
                        !!option.label?.toLowerCase().startsWith(searchQuery.toLowerCase())
                      }
                      value={stores.filter((store: any) => matchingContent?.storeID?.includes(store?.id))}
                      onChange={(selectedStores: SelectableValue<string>) =>
                        onUpdateContent(matchingContent, ReportContentKey.STORE_ID, selectedStores)
                      }
                      options={stores.map((store: any) => ({ label: store.name, value: store }))}
                    />
                  </Field>
                  <Field label={intl.get('lookback_period')} description={intl.get('lookback_period_description')}>
                    <Select
                      value={getLookbacks().filter((lookback: any) => lookback.value === matchingContent?.lookback)}
                      options={getLookbacks()}
                      onChange={(selectedLookback: SelectableValue) => {
                        onUpdateContent(matchingContent, ReportContentKey.LOOKBACK, selectedLookback);
                      }}
                    />
                  </Field>
                </>
              )}
            </div>
          </li>
        );
      })}
    </ol>
  );
};
