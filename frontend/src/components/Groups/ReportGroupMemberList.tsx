import React, { FC, useState } from 'react';
import classNames from 'classnames';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Checkbox, Field, Icon, Input, Label, Legend, Spinner, Tooltip } from '@grafana/ui';
import { useQuery } from 'react-query';
import { getUsers, getGroupMembers, deleteReportGroupMembership, createReportGroupMembership } from 'api';
import { CreateGroupMemberVars, ReportGroup, ReportGroupMember, User } from 'common/types';

import { useOptimisticMutation } from 'hooks/useOptimisticMutation';
import { useWindowSize } from 'hooks/useWindowResize';
import { useDatasourceID } from 'hooks';

type Props = {
  reportGroup: ReportGroup | null;
};

const listStyle = classNames({
  'card-section': true,
  'card-list-layout-grid': true,
  'card-list-layout-list': true,
  'card-list': true,
});

const marginForCheckbox = css`
  margin-right: 10px;
`;

export const ReportGroupMemberList: FC<Props> = ({ reportGroup }) => {
  const datasourceID = useDatasourceID();
  console.log(datasourceID);
  const { id: reportGroupID } = reportGroup ?? {};
  const { height } = useWindowSize();

  const [searchTerm, setSearchTerm] = useState<string>('');

  const { data: users, error } = useQuery<User[], Error>(`users-${datasourceID}`, () => getUsers(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    retry: 0,
  });

  console.log(error);

  const { data: groupMembers } = useQuery<ReportGroupMember[]>(['groupMembers', reportGroupID], getGroupMembers);

  const [createMembership] = useOptimisticMutation<
    ReportGroupMember[],
    ReportGroupMember,
    CreateGroupMemberVars,
    ReportGroupMember[]
  >(
    ['groupMembers', reportGroupID],
    createReportGroupMembership,
    variables => ({ ...variables, userID: variables.user.id, id: '' }),
    (prevState, optimisticValue) => {
      if (prevState) {
        return [...prevState, optimisticValue];
      } else {
        return prevState;
      }
    },
    []
  );

  const [deleteMembership] = useOptimisticMutation<
    ReportGroupMember[],
    ReportGroupMember,
    ReportGroupMember,
    ReportGroupMember[]
  >(
    ['groupMembers', reportGroupID],
    deleteReportGroupMembership,
    groupMember => groupMember,
    (prevState, optimisticValue) => {
      if (prevState) {
        return prevState.filter(member => member.id !== optimisticValue.id);
      } else {
        return prevState;
      }
    },
    []
  );

  const onToggleMember = (user: User) => {
    const { id: userID } = user;
    const exists = groupMembers?.find((reportMember: ReportGroupMember) => reportMember.userID === userID);

    if (exists) {
      deleteMembership(exists);
    } else {
      createMembership({ user, reportGroupID: reportGroupID ?? '' });
    }
  };

  return (
    <div style={{}}>
      <div style={{ marginTop: '25px', display: 'flex', alignItems: 'center' }}>
        <Tooltip placement="top" content={intl.get('users_tooltip')} theme={'info'}>
          <Icon
            name="info-circle"
            size="sm"
            style={{ marginLeft: '10px', marginRight: '10px', marginBottom: '16px' }}
          />
        </Tooltip>
        <Legend>
          {error ? 'There was a problem finding users. Is your Datasource setup correctly?' : intl.get('users')}
        </Legend>
      </div>
      {!users ? (
        !error && <Spinner style={{ flex: 1, display: 'flex', justifyContent: 'center' }} />
      ) : (
        <>
          <Field label="Search for users">
            <Input
              name="search"
              css=""
              placeholder="Search for the user"
              type="text"
              prefix={<Icon name="search" />}
              suffix={<Icon name="trash-alt" onClick={() => setSearchTerm('')} />}
              onChange={e => {
                const { value } = e.target as HTMLInputElement;
                return setSearchTerm(value);
              }}
            />
          </Field>
          <Field>
            <Label description="Option description">{searchTerm}</Label>
          </Field>

          <ol className={listStyle} style={{ maxHeight: `${(height ?? 0) / 2}px`, overflow: 'scroll' }}>
            {users
              .filter(user => {
                const match = user.name.toString().toLowerCase().indexOf(searchTerm.toLowerCase()) > -1;
                return match;
              })
              .map((user: User) => {
                const { name, e_mail, id } = user;
                const isChecked = groupMembers?.find((groupMember: ReportGroupMember) => groupMember.userID === id);
                return (
                  <li className="card-item-wrapper" style={e_mail ? { cursor: 'pointer' } : {}} key={id}>
                    <div className={'card-item'} onClick={() => e_mail && onToggleMember(user)}>
                      <div className="card-item-body">
                        <div className={marginForCheckbox}>
                          <Checkbox value={!!isChecked} css="" />
                        </div>
                        <div className="card-item-details">
                          <div className="card-item-name">{name}</div>
                          <div className="card-item-type">{e_mail ? e_mail : intl.get('no_email')}</div>
                        </div>
                      </div>
                    </div>
                  </li>
                );
              })}
          </ol>
        </>
      )}
    </div>
  );
};
