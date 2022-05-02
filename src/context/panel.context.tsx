import { getPanels } from 'api';
import { useDatasourceID } from 'hooks';
import React, { useState, useEffect } from 'react';
import { useQuery } from 'react-query';
import { Panel } from 'types';
import { panelUsesMacro } from 'utils';

const PanelContext = React.createContext<any | null>(null);

const PanelProvider: React.FC = ({ children }) => {
  const [panelDetails, setPanelDetails] = useState<any[]>([]);

  const datasourceID = useDatasourceID();

  const { data: panels } = useQuery<Panel[], Error>(['panels'], () => getPanels(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  useEffect(() => {
    if (panels) {
      const newPanelDetails = panels.map((panel) => {
        const usesMacro = panelUsesMacro(panel.rawSql);
        const usesVariables = panel.variables.length > 0;

        return {
          id: '',
          scheduleID: '',
          panelID: panel.id,
          dashboardID: panel.dashboardID,
          lookback: usesMacro ? 1 : 0,
          variables: usesVariables ? '' : null,
        };
      });

      setPanelDetails(newPanelDetails);
    }
  }, [panels]);

  return (
    <PanelContext.Provider
      value={{
        panels,
        panelDetails,
        setPanelDetails,
      }}
    >
      {children}
    </PanelContext.Provider>
  );
};

export { PanelContext, PanelProvider };
