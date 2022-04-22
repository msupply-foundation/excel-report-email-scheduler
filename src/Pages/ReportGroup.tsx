import React, { useState } from 'react';
import { css, cx } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { Button, Card, ConfirmModal, HorizontalGroup, LinkButton, Spinner, Tag, useStyles2 } from '@grafana/ui';
import intl from 'react-intl-universal';
import { EmptyListCTA } from 'components/common';
import { prefixRoute } from '../utils';
import { PLUGIN_BASE_URL, ROUTES } from '../constants';
import { ReportGroupType } from 'types';
import { useMutation, useQuery } from 'react-query';
import { deleteReportGroup, getReportGroups } from 'api/ReportGroup';
import { useToggle } from 'hooks';

const EmptyList = () => {
  return (
    <EmptyListCTA
      title="You haven't created any report groups yet."
      buttonTitle={'Create new report group'}
      buttonIcon="users-alt"
      buttonLink={`${prefixRoute(ROUTES.REPORT_GROUP)}/create`}
      proTip="Report groups are groups containing users to send the report to."
      proTipLink=""
      proTipLinkTitle=""
      proTipTarget="_blank"
    />
  );
};

const ReportGroup = () => {
  const styles = useStyles2(getStyles);

  const [deleteAlertIsOpen, setDeleteAlertIsOpen] = useToggle(false);
  const [deleteReportGroupID, setDeleteReportGroupID] = useState('');

  const { mutate: deleteGroup } = useMutation(deleteReportGroup, {
    onSuccess: () => {
      console.log('deleted success');
      refetchReportGroups();
      return;
    },
  });

  const onConfirmDeleteGroup = () => {
    deleteGroup(deleteReportGroupID);
    setDeleteAlertIsOpen();
  };

  const {
    data: reportGroups,
    isLoading,
    refetch: refetchReportGroups,
  } = useQuery<ReportGroupType[], Error>(`reportGroups`, getReportGroups, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const onReportGroupDelete = (reportGroupID: string) => {
    setDeleteAlertIsOpen();
    setDeleteReportGroupID(reportGroupID);
  };

  if (isLoading) {
    return <Spinner />;
  }

  if (!!!reportGroups || !reportGroups.length) {
    return <EmptyList />;
  }

  return (
    <div>
      <div className={styles.adjustButtonToRight}>
        <LinkButton icon="plus-circle" key="create" variant="primary" href={`${PLUGIN_BASE_URL}/report-groups/create`}>
          {intl.get('add_report_group')}
        </LinkButton>
      </div>
      <ul className={styles.list}>
        {reportGroups.map((reportGroup) => {
          return (
            <li key={reportGroup.id}>
              <Card
                className={cx(styles.card, 'card-parent')}
                href={`${PLUGIN_BASE_URL}/report-groups/edit/${reportGroup.id}`}
              >
                <Card.Heading className={styles.heading}>{reportGroup.name}</Card.Heading>
                <Card.Description className={styles.description}>{reportGroup.description}</Card.Description>
                {reportGroup.members && (
                  <Card.Meta>
                    {[
                      <HorizontalGroup
                        key="members"
                        wrap={true}
                        style={{ marginBottom: '25px' }}
                        align="flex-start"
                        justify="flex-start"
                      >
                        {reportGroup.members.map(({ id, name, email }: any) => {
                          return <Tag key={id} icon="user" name={`${name} <${email} >`} />;
                        })}
                      </HorizontalGroup>,
                    ]}
                  </Card.Meta>
                )}
                <Card.Actions className={styles.actions}>
                  <LinkButton
                    icon="cog"
                    key="edit"
                    variant="secondary"
                    href={`${PLUGIN_BASE_URL}/report-groups/edit/${reportGroup.id}`}
                  >
                    Edit
                  </LinkButton>
                  <Button
                    key="delete"
                    icon="trash-alt"
                    variant="destructive"
                    onClick={(e) => onReportGroupDelete(reportGroup.id)}
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

const getStyles = (theme: GrafanaTheme2) => ({
  marginTop: css`
    margin-top: ${theme.spacing(2)};
  `,
  list: css({
    listStyle: 'none',
    display: 'grid',
    // gap: '8px', Add back when legacy support for old Card interface is dropped
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

export { ReportGroup, EmptyList };
