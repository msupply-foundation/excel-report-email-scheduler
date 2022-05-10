import React, { useState } from 'react';
import { css, cx } from '@emotion/css';
import intl from 'react-intl-universal';
import { GrafanaTheme2 } from '@grafana/data';
import { getLookbacks, PLUGIN_BASE_URL, ROUTES } from '../../constants';
import { prefixRoute } from '../../utils/navigation';
import { EmptyListCTA, Loading } from 'components/common';
import { ScheduleType } from 'types';
import { useQuery, useMutation } from 'react-query';
import { deleteSchedule, getSchedules } from 'api';
import { Button, Card, ConfirmModal, HorizontalGroup, LinkButton, Spinner, Tag, useStyles2 } from '@grafana/ui';
import { useToggle } from 'hooks';
import { PanelContext, PanelProvider } from 'context';

const Schedule: React.FC = () => {
  const styles = useStyles2(getStyles);

  const [deleteAlertIsOpen, setDeleteAlertIsOpen] = useToggle(false);
  const [deleteScheduleID, setDeleteScheduleID] = useState('');

  const {
    data: schedules,
    isLoading,
    refetch: refetchSchedules,
  } = useQuery<ScheduleType[], Error>(`scheduled`, getSchedules, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const onScheduleDelete = (reportGroupID: string) => {
    setDeleteAlertIsOpen();
    setDeleteScheduleID(reportGroupID);
  };

  const { mutate: deleteScheduleMutate } = useMutation(deleteSchedule, {
    onSuccess: () => {
      refetchSchedules();
      return;
    },
  });

  const onConfirmDeleteGroup = () => {
    deleteScheduleMutate(deleteScheduleID);
    setDeleteAlertIsOpen();
  };

  if (isLoading) {
    return <Spinner />;
  }

  if (!!!schedules || !schedules.length) {
    return <EmptyList />;
  }

  return (
    <div>
      <div className={styles.adjustButtonToRight}>
        <LinkButton icon="plus-circle" key="create" variant="primary" href={`${PLUGIN_BASE_URL}/schedules/create`}>
          {intl.get('add_schedule')}
        </LinkButton>
      </div>
      <ul className={styles.list}>
        {schedules.map((schedule) => {
          return (
            <li key={schedule.id}>
              <Card
                className={cx(styles.card, 'card-parent')}
                href={`${PLUGIN_BASE_URL}/schedules/edit/${schedule.id}`}
              >
                <Card.Heading className={styles.heading}>{schedule.name}</Card.Heading>
                <Card.Description className={styles.description}>{schedule.description}</Card.Description>
                {schedule.panelDetails && (
                  <Card.Meta>
                    {
                      <HorizontalGroup
                        spacing="lg"
                        key="panelDetails"
                        wrap={true}
                        style={{ marginBottom: '25px' }}
                        align="flex-start"
                        justify="flex-start"
                      >
                        <div style={{ flexWrap: 'wrap', lineHeight: '2em' }}>
                          <PanelProvider>
                            <PanelContext.Consumer>
                              {({ panels }) => {
                                const lookbacks = getLookbacks();

                                if (!panels || !panels.length) {
                                  return <Loading />;
                                }

                                return (
                                  !!panels &&
                                  schedule.panelDetails.map(({ id, panelID, lookback, variables }: any) => {
                                    const panel = panels.find((panel: any) => panel.id === panelID);
                                    const lookbackLabel = lookbacks.find((el) => el.value === lookback)?.label;

                                    if (!panel) {
                                      return false;
                                    }

                                    return (
                                      <Tag
                                        key={id}
                                        icon="user"
                                        className={styles.tag}
                                        name={
                                          `Panel: ${panel.title}` +
                                          (!!lookback ? ` | Lookback: ${lookbackLabel}` : '') +
                                          (!!variables ? ` | Variables: ${variables}` : '')
                                        }
                                      />
                                    );
                                  })
                                );
                              }}
                            </PanelContext.Consumer>
                          </PanelProvider>
                        </div>
                      </HorizontalGroup>
                    }
                  </Card.Meta>
                )}
                <Card.Actions className={styles.actions}>
                  <LinkButton
                    icon="cog"
                    key="edit"
                    variant="secondary"
                    href={`${PLUGIN_BASE_URL}/schedules/edit/${schedule.id}`}
                  >
                    Edit
                  </LinkButton>
                  <Button
                    key="delete"
                    icon="trash-alt"
                    variant="destructive"
                    onClick={(e) => onScheduleDelete(schedule.id)}
                  >
                    {intl.get('delete')}
                  </Button>
                </Card.Actions>
              </Card>
            </li>
          );
        })}
      </ul>
      <ConfirmModal
        isOpen={deleteAlertIsOpen}
        title={intl.get('delete_report_group')}
        body={intl.get('delete_report_group_question')}
        confirmText={intl.get('delete')}
        icon="exclamation-triangle"
        onConfirm={onConfirmDeleteGroup}
        onDismiss={setDeleteAlertIsOpen}
      />
    </div>
  );
};

const EmptyList = () => {
  return (
    <EmptyListCTA
      title="You haven't created any schedule(s) yet."
      buttonTitle={'Create new schedule'}
      buttonIcon="calendar-alt"
      buttonLink={`${prefixRoute(ROUTES.SCHEDULES)}/create`}
      proTip="Schedules are set with report groups and panels. They run on fixed schedules recursively sending excel reports through email."
    />
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  marginTop: css`
    margin-top: ${theme.spacing(2)};
  `,
  tag: css`
    margin-right: 5px;
    padding: 5px;
  `,
  list: css({
    listStyle: 'none',
    display: 'grid',
    gap: '8px',
  }),
  adjustButtonToRight: css`
    display: flex;
    justify-content: flex-end;
    margin-bottom: 10px;
  `,
  heading: css({
    fontSize: theme.v1.typography.heading.h5,
    fontWeight: 'inherit',
  }),
  card: css({
    gridTemplateAreas: `
      "Figure   Heading   Actions"
      "Figure Description Actions"
      "Figure    Meta     Actions"
      "Figure     -       Actions"`,
  }),
  description: css({
    margin: '0px',
    fontSize: theme.typography.size.sm,
  }),
  actions: css({
    position: 'relative',
    alignSelf: 'center',
    marginTop: '0px',
    opacity: 0,

    '.card-parent:hover &, .card-parent:focus-within &': {
      opacity: 1,
    },
  }),
});

export { Schedule };
