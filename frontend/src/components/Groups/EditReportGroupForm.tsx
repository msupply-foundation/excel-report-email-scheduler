import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Icon, InlineFormLabel, Input, Legend, Tooltip } from '@grafana/ui';
import { ReportGroupKey } from 'common/enums';
import { ReportGroup } from 'common/types';

const container = css`
  display: flex;
  flex-direction: column;
  flex: 1;
  padding-right: 30px;
`;

type Props = {
  onUpdate: (key: ReportGroupKey, value: string) => void;
  reportGroup: ReportGroup | null;
};

const flexWrapping = css`
  display: flex;
  flex-direction: row;
  flex: 1;
  flex-basis: calc(50% - 20px);
`;

export const EditReportGroupForm: FC<Props> = ({ onUpdate, reportGroup }) => (
  <div className={container}>
    <div style={{ marginTop: '25px', display: 'flex', alignItems: 'center' }}>
      <Tooltip placement="top" content={intl.get('edit_details_report_group_tooltip')} theme={'info'}>
        <Icon name="info-circle" size="sm" style={{ marginLeft: '10px', marginRight: '10px', marginBottom: '16px' }} />
      </Tooltip>
      <Legend>{intl.get('edit_details')}</Legend>
    </div>
    <div className={flexWrapping}>
      <InlineFormLabel tooltip={intl.get('group_name')}>{intl.get('name')}</InlineFormLabel>

      <Input
        onChange={({ currentTarget: { value } }) => onUpdate(ReportGroupKey.NAME, value)}
        name={intl.get('name')}
        defaultValue={reportGroup?.name}
        placeholder={intl.get('name')}
        css=""
      />
    </div>

    <div className={flexWrapping}>
      <InlineFormLabel tooltip={intl.get('description')}>{intl.get('group_description')}</InlineFormLabel>

      <Input
        onChange={({ currentTarget: { value } }) => onUpdate(ReportGroupKey.DESCRIPTION, value)}
        name={intl.get('description')}
        defaultValue={reportGroup?.description}
        placeholder={intl.get('description')}
        css=""
      />
    </div>
  </div>
);
