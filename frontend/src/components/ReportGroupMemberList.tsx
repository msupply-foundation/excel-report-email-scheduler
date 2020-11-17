import React, { FC } from 'react';
import classNames from 'classnames';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Checkbox } from '@grafana/ui';
import { queryCache, useMutation, useQuery } from 'react-query';
import { getUsers, getGroupMembers, deleteReportGroupMembership, createReportGroupMembership } from 'api';
import { ReportGroupMember, User } from 'common/types';
import { ReportGroup } from './ReportSchedulesTab';

type Props = {
  reportGroup: ReportGroup | null;
  datasourceID: number;
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

export const ReportGroupMemberList: FC<Props> = ({ reportGroup, datasourceID }) => {
  const { id: reportGroupID } = reportGroup ?? {};

  const { data: users } = useQuery<User[]>('users', () => getUsers(datasourceID));
  const { data: groupMembers } = useQuery<ReportGroupMember[]>(['groupMembers', reportGroupID], getGroupMembers);

  const [createMembership] = useMutation(createReportGroupMembership, {
    onSuccess: () => queryCache.refetchQueries(['groupMembers', reportGroupID]),
  });
  const [deleteMembership] = useMutation(deleteReportGroupMembership, {
    onSuccess: () => queryCache.refetchQueries(['groupMembers', reportGroupID]),
  });

  const onToggleMember = (user: User) => {
    const { id: userID } = user;
    const exists = groupMembers?.find((reportMember: ReportGroupMember) => reportMember.userID === userID);

    if (exists) {
      deleteMembership(exists);
    } else {
      createMembership({ user, reportGroupID });
    }
  };

  return (
    <ol className={listStyle}>
      {users?.map((user: User) => {
        console.log('groupMembers', groupMembers);
        const { name, e_mail, id } = user;
        const isChecked = groupMembers?.find((groupMember: ReportGroupMember) => groupMember.userID === id);
        return (
          <li className="card-item-wrapper" style={{ cursor: 'pointer' }}>
            <div className={'card-item'} onClick={() => onToggleMember(user)}>
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
  );
};
