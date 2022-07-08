import { SelectableValue } from '@grafana/data';
import { getPanels } from 'api';
import { useDatasourceID } from 'hooks';
import React, { useState, useEffect } from 'react';
import { useQuery } from 'react-query';
import { ContentVariables, Panel, PanelDetails, SelectableVariable } from 'types';
import { parseOrDefault } from 'utils';

type PanelContextProps = {
  panels: Panel[];
  panelDetails: PanelDetails[];
  setPanelDetails: any;
  onUpdateLookback: (content: PanelDetails) => (selectableValue: SelectableValue) => void;
  onUpdateVariable: (
    content: PanelDetails,
    panel: Panel
  ) => (variableName: string) => (selectedValue: SelectableValue) => void;
};

const panelContextDefault = {
  panels: [],
  panelDetails: [],
  setPanelDetails: (panelDetails: PanelDetails[]) => {},
  onUpdateVariable:
    (content: PanelDetails, panel: Panel) => (variableName: string) => (selectedValue: SelectableValue) => {},
  onUpdateLookback: (content: PanelDetails) => (selectableValue: SelectableValue) => {},
};

const PanelContext = React.createContext<PanelContextProps>(panelContextDefault);

const PanelProvider: React.FC = ({ children }) => {
  const [panelDetails, setPanelDetails] = useState<PanelDetails[] | []>([]);

  const datasourceID = useDatasourceID();

  const { data: panels } = useQuery<Panel[], Error>('all-panels', () => getPanels(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  useEffect(() => {
    if (panels) {
      const newPanelDetails = panels.map((panel) => ({
        id: '',
        scheduleID: '',
        panelID: panel.id,
        dashboardID: panel.dashboardID,
        lookback: '',
        variables: '',
      }));

      setPanelDetails(newPanelDetails);
    }
  }, [panels]);

  const onUpdateLookback = (content: PanelDetails) => (selectableValue: SelectableValue) => {
    setPanelDetails((prevPanels: any) => {
      const myIndex = prevPanels.findIndex(
        (el: any) => el.panelID === content.panelID && el.dashboardID === content.dashboardID
      );

      return [
        ...prevPanels.slice(0, myIndex),
        { ...prevPanels[myIndex], lookback: selectableValue.value },
        ...prevPanels.slice(myIndex + 1),
      ];
    });
  };

  const onUpdateVariable =
    (content: PanelDetails, panel: Panel) =>
    (variableName: string) =>
    (selectableValue: SelectableValue<SelectableVariable[]>) => {
      const newVariable = selectableValue.map(({ value }: SelectableValue) => value.value);
      const newVariables = parseOrDefault<ContentVariables>(content.variables, {});

      newVariables[variableName] = newVariable;

      setPanelDetails((prevPanels: any) => {
        const myIndex = prevPanels.findIndex(
          (el: any) => el.panelID === panel.id && el.dashboardID === panel.dashboardID
        );

        return [
          ...prevPanels.slice(0, myIndex),
          { ...prevPanels[myIndex], variables: JSON.stringify(newVariables) },
          ...prevPanels.slice(myIndex + 1),
        ];
      });
    };

  return (
    <PanelContext.Provider
      value={{
        panels: panels || [],
        panelDetails,
        setPanelDetails,
        onUpdateLookback,
        onUpdateVariable,
      }}
    >
      {children}
    </PanelContext.Provider>
  );
};

export { PanelContext, PanelProvider, PanelContextProps };
