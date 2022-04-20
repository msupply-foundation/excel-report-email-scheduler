import React from 'react';
import { css, cx } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { Button, Card, HorizontalGroup, Spinner, Tag, useStyles2 } from '@grafana/ui';
import intl from 'react-intl-universal';
import { EmptyListCTA } from 'components/common';
import { prefixRoute } from '../utils';
import { PLUGIN_BASE_URL, ROUTES } from '../constants';
import { ReportGroupType } from 'types';
import { useQuery } from 'react-query';
import { getReportGroups } from 'api/ReportGroup';

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

  const { data: reportGroups, isLoading } = useQuery<ReportGroupType[], Error>(`reportGroups`, getReportGroups, {
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  if (isLoading) {
    return <Spinner />;
  }

  if (reportGroups === undefined || reportGroups.length <= 0) {
    return <EmptyList />;
  }

  return (
    <div>
      <div className={styles.adjustButtonToRight}>
        <Button onClick={() => {}} variant="primary">
          {intl.get('add_report_group')}
        </Button>
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
                  <Button icon="cog" key="edit" variant="secondary">
                    Edit
                  </Button>
                  <Button key="edit">Delete</Button>
                </Card.Actions>
              </Card>
            </li>
          );
        })}
      </ul>
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
