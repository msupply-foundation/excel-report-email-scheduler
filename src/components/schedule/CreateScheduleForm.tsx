import { Form, FormAPI, LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { createSchedule, getReportGroups, getScheduleByID } from 'api';
import { CreateScheduleFormPartial } from 'components';
import { PLUGIN_BASE_URL } from '../../constants';
import { PanelContext } from 'context';
import React, { useContext } from 'react';
import { useMutation, useQuery } from 'react-query';
import { ScheduleType, ReportGroupType, PanelDetails } from 'types';
import { useHistory, useParams } from 'react-router-dom';
import { css } from '@emotion/css';

const defaultFormValues: ScheduleType = {
  id: '',
  name: '',
  description: '',
  interval: 0,
  time: '',
  reportGroupID: '',
  day: 1,
  panels: [],
  panelDetails: [],
};

const CreateScheduleForm: React.FC = () => {
  const { panelDetails } = useContext(PanelContext);
  const style = useStyles2(getStyles);

  const [defaultSchedule, setDefaultSchedule] = React.useState<ScheduleType>(defaultFormValues);

  const history = useHistory();

  const { id: scheduleIdToEdit } = useParams<{ id: string }>();
  const isEditMode = !!scheduleIdToEdit;
  const [ready, setReady] = React.useState(false);

  const {
    data: defaultScheduleFetched,
    isLoading: isScheduleByIDLoading,
    isRefetching,
  } = useQuery<ScheduleType, Error>(`schedules-${scheduleIdToEdit}`, () => getScheduleByID(scheduleIdToEdit), {
    enabled: isEditMode,
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    onError: () => {
      history.push(`${PLUGIN_BASE_URL}/schedules/`);
      return;
    },
  });

  React.useEffect(() => {
    if (!!defaultScheduleFetched) {
      setDefaultSchedule({
        ...defaultScheduleFetched,
      });
    }
  }, [defaultScheduleFetched]);

  React.useEffect(() => {
    if (!isEditMode) {
      setReady(true);
    }

    if (isEditMode && !isRefetching && !isScheduleByIDLoading) {
      setReady(true);
    }
  }, [isEditMode, isRefetching, isScheduleByIDLoading]);

  const { data: reportGroups } = useQuery<ReportGroupType[], Error>(`reportGroups`, getReportGroups, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const createScheduleMutation = useMutation((newSchedule: ScheduleType) => createSchedule(newSchedule), {
    onSuccess: () => {
      history.push(`${PLUGIN_BASE_URL}/schedules/`);
      return;
    },
  });

  const submitCreateSchedule = (data: ScheduleType) => {
    console.log(data);
    const selectedPanels = panelDetails.filter((detail: PanelDetails) => data.panels.includes(detail.panelID));
    data.panelDetails = selectedPanels;

    console.log(data);
    createScheduleMutation.mutate(data);
  };

  if (!ready) {
    return (
      <div className={style.loadingWrapper}>
        <LoadingPlaceholder text="Loading..." />
      </div>
    );
  }

  return (
    <Form onSubmit={submitCreateSchedule} validateOnMount={false} defaultValues={defaultSchedule} validateOn="onSubmit">
      {(props: FormAPI<ScheduleType>) => {
        return (
          <CreateScheduleFormPartial
            isEditMode={isEditMode}
            defaultSchedule={defaultSchedule}
            {...props}
            reportGroups={reportGroups}
          />
        );
      }}
    </Form>
  );
};

const getStyles = () => ({
  loadingWrapper: css`
    display: flex;
    height: 50vh;
    align-items: center;
    justify-content: center;
  `,
});

export { CreateScheduleForm };
