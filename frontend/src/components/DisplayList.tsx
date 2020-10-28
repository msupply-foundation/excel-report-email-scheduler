import React, { FC } from 'react';
import classNames from 'classnames';
import { css } from 'emotion';
import { AsyncMultiSelect, Checkbox, MultiSelect, Select } from '@grafana/ui';
import { getStores } from '../api';
import { SelectableValue } from '@grafana/data';

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

type Props = {
  onRowPress: (toggle: any) => void;
  data: any;
  titleKey: string;
  descriptionKey?: string;
  withChecks?: boolean;
  checked?: any;
};

export const DisplayList: FC<Props> = ({ titleKey, checked, descriptionKey, onRowPress, data = [], withChecks }) => {
  const listStyle = classNames({
    'card-section': true,
    'card-list-layout-grid': true,
    'card-list-layout-list': true,
  });

  return (
    <section className={listStyle}>
      <ol className="card-list">
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
                </div>
              </div>
            </li>
          );
        })}
      </ol>
    </section>
  );
};
