import React, { FC } from 'react';
import classNames from 'classnames';
import { css } from 'emotion';
import { AsyncMultiSelect, Checkbox, MultiSelect, Select } from '@grafana/ui';
import { getStores, updateReportContent } from '../api';
import { SelectableValue } from '@grafana/data';
import { queryCache, useMutation, useQuery } from 'react-query';

const lookback = [
  { label: '1 Day', value: 24 * 60 * 1000 * 60 },
  { label: '2 Days', value: 2 * 24 * 60 * 1000 * 60 },
  { label: '3 Days', value: 3 * 24 * 60 * 1000 * 60 },
  { label: '1 Week', value: 7 * 24 * 60 * 1000 * 60 },
  { label: '2 Weeks', value: 14 * 24 * 60 * 1000 * 60 },
  { label: '4 Weeks', value: 28 * 24 * 60 * 1000 * 60 },
  { label: '3 Months', value: 28 * 3 * 24 * 60 * 1000 * 60 },
  { label: '6 Months', value: 28 * 6 * 24 * 60 * 1000 * 60 },
  { label: '1 Year', value: 356 * 24 * 60 * 1000 * 60 },
];

const listStyle = classNames({
  'card-section': true,
  'card-list-layout-grid': true,
  'card-list-layout-list': true,
  'card-list': true,
});

type Props = {
  scheduleID: string | any;
  onRowPress: (toggle: any) => void;
  data: any;
  titleKey: string;
  descriptionKey?: string;
  withChecks?: boolean;
  checked?: any;
};

export const PanelList: FC<Props> = ({
  titleKey,
  checked,
  descriptionKey,
  scheduleID,
  onRowPress,
  data = [],
  withChecks,
}) => {
  const [updateContent] = useMutation(updateReportContent, {
    onSuccess: () => queryCache.refetchQueries(['reportContent', scheduleID]),
  });
  const { data: stores } = useQuery('stores', getStores);

  const onUpdateContent = (panel: any, key: string, value: any) => {
    let newValue = value;
    const content = checked[panel?.id];
    console.log('newValue', value);
    if (key === 'storeID') {
      newValue = value.map((store: any) => store?.value?.id).join(', ');
    }

    const newState = { ...content, [key]: newValue };

    updateContent(newState);
  };

  return (
    <ol className={listStyle}>
      {data?.map((datum: any) => {
        console.log(datum);
        return (
          <li className="card-item-wrapper" style={{ cursor: 'pointer' }}>
            <div
              className={css`
                display: flex;
                flex-direction: row;
              `}
            >
              <div className={'card-item'} style={{ flex: 4 }}>
                <div className="card-item-body" onClick={() => onRowPress(datum)}>
                  {withChecks && (
                    <div
                      className={css`
                        margin-right: 10px;
                      `}
                    >
                      <Checkbox value={!!checked?.[datum.id]} css="" />
                    </div>
                  )}
                  <div className="card-item-details">
                    <div className="card-item-name">{datum[titleKey]}</div>
                    {descriptionKey && <div className="card-item-type">{datum[descriptionKey]}</div>}
                  </div>
                </div>

                {!!checked?.[datum.id] && (
                  <MultiSelect
                    closeMenuOnSelect={false}
                    filterOption={(option: SelectableValue, searchQuery: string) =>
                      !!option.label?.toLowerCase().startsWith(searchQuery.toLowerCase())
                    }
                    value={stores.filter((store: any) => {
                      return checked?.[datum.id]?.storeID?.includes(store?.id);
                    })}
                    onChange={e => {
                      onUpdateContent(datum, 'storeID', e);
                    }}
                    options={stores.map((store: any) => ({ label: store.name, value: store }))}
                  />
                )}
              </div>
            </div>
          </li>
        );
      })}
    </ol>
  );
};
