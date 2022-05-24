import { Form, FormAPI, PageToolbar, Spinner, ToolbarButton } from '@grafana/ui';
import { createSchedule, getReportGroups, getScheduleByID, sendTestEmail } from 'api';
import { CreateScheduleFormPartial, Loading } from 'components';
import { PLUGIN_BASE_URL, ROUTES } from '../../constants';
import { PanelContext } from 'context';
import React, { useContext } from 'react';
import { useMutation, useQuery } from 'react-query';
import { ScheduleType, ReportGroupType, PanelDetails } from 'types';
import { useHistory, useParams } from 'react-router-dom';
import { prefixRoute } from 'utils';
import intl from 'react-intl-universal';

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
  const { panelDetails, setPanelDetails } = useContext(PanelContext);

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

      setPanelDetails((prevDetails: PanelDetails[]) =>
        prevDetails.map(
          (prevDetail: PanelDetails) =>
            defaultScheduleFetched.panelDetails.find(
              (defaultDetail) =>
                defaultDetail.panelID === prevDetail.panelID && defaultDetail.dashboardID === prevDetail.dashboardID
            ) || prevDetail
        )
      );
    }
  }, [defaultScheduleFetched, setPanelDetails]);

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

  const { mutate: testEmails, isLoading: isSendTestEmailLoading } = useMutation(sendTestEmail);

  const submitCreateSchedule = (data: ScheduleType) => {
    const selectedPanels = panelDetails.filter((detail: PanelDetails) =>
      data.panels.find((panel) => panel.panelID === detail.panelID && panel.dashboardID === detail.dashboardID)
    );
    data.panelDetails = selectedPanels;
    createScheduleMutation.mutate(data);
  };

  if (!ready) {
    return <Loading />;
  }

  return (
    <>
      <PageToolbar
        parent="Schedules"
        titleHref="#"
        parentHref={prefixRoute(ROUTES.SCHEDULES)}
        title={`${isEditMode ? 'Edit "' + defaultSchedule?.name + '"' : 'New'}`}
        onGoBack={() => history.push(prefixRoute(ROUTES.SCHEDULES))}
      >
        {isEditMode &&
          (isSendTestEmailLoading ? (
            <Spinner />
          ) : (
            <ToolbarButton icon="envelope" onClick={() => testEmails(defaultSchedule.id)}>
              {intl.get('send_test_emails')}
            </ToolbarButton>
          ))}
      </PageToolbar>
      <Form
        style={{ marginTop: '30px' }}
        onSubmit={submitCreateSchedule}
        validateOnMount={false}
        defaultValues={defaultSchedule}
        validateOn="onSubmit"
      >
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
    </>
  );
};

export { CreateScheduleForm };
