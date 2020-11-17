import { ContentVariables, Panel, ReportContent, Store } from 'common/types';
import React, { FC } from 'react';

import { css } from 'emotion';
import { Checkbox } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { ReportContentKey } from 'common/enums';
import { useOptimisticMutation } from 'hooks';
import { updateReportContent } from 'api';

import { parseOrDefault } from 'common';
import { PanelVariables } from './PanelVariables';

type Props = {
  panel: Panel;
  reportContent: ReportContent | null;
  onToggle: (panel: Panel) => Promise<void>;
  stores?: Store[];
  scheduleID: string;
};

const marginForCheckbox = css`
  margin-right: 10px;
`;

export const PanelItem: FC<Props> = ({ panel, reportContent, onToggle, scheduleID }) => {
  const { title, description } = panel;

  const [updateContent] = useOptimisticMutation<ReportContent[], ReportContent, ReportContent, ReportContent[]>(
    ['reportContent', scheduleID],
    updateReportContent,
    (content: ReportContent): ReportContent => content,
    (prevState: ReportContent[] | undefined, optimisticValue: ReportContent) => {
      const { id: optimisticID } = optimisticValue;
      if (prevState) {
        const idx = prevState.findIndex(({ id }) => id === optimisticID);
        prevState[idx] = optimisticValue;
        return [...prevState];
      } else {
        return prevState;
      }
    },
    []
  );

  const onUpdateContent = (content: ReportContent) => (
    key: ReportContentKey,
    selectableValue: SelectableValue<String | Number | Store>
  ) => {
    let newValue = selectableValue.value;
    if (key === ReportContentKey.STORE_ID && Array.isArray(selectableValue)) {
      newValue = selectableValue.map((selected: SelectableValue<Store>) => selected.value?.id).join(', ');
    }

    const newState = { ...content, [key]: newValue };
    updateContent(newState);
  };

  const onUpdateVariable = (content: ReportContent) => (variableName: string) => (selectableValue: SelectableValue) => {
    const newVariable = selectableValue.map(({ value }: SelectableValue) => value.value);
    const newVariables = parseOrDefault<ContentVariables>(content.variables, {});
    newVariables[variableName] = newVariable;
    updateContent({ ...content, variables: JSON.stringify(newVariables) });
  };

  return (
    <li className="card-item-wrapper" style={{ cursor: 'pointer' }}>
      <div className={'card-item'}>
        <div className="card-item-body" onClick={() => onToggle(panel)}>
          <div className={marginForCheckbox}>
            <Checkbox value={!!reportContent} css="" />
          </div>

          <div className="card-item-details">
            <div className="card-item-name">{title}</div>
            <div className="card-item-type">{description}</div>
          </div>
        </div>

        {reportContent && (
          <PanelVariables
            storeIDs={reportContent.storeID}
            panel={panel}
            lookback={reportContent.lookback}
            variables={reportContent.variables}
            onUpdateContent={onUpdateContent(reportContent)}
            onUpdateVariable={onUpdateVariable(reportContent)}
          />
        )}
      </div>
    </li>
  );
};
