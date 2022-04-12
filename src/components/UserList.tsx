import React, { useState, useCallback, useEffect } from 'react';
import {
  InlineField,
  Input,
  Icon,
  InlineSwitch,
  VerticalGroup,
  HorizontalGroup,
  Pagination,
  Checkbox,
  Tag,
  FieldSet,
  EmptySearchResult,
  useStyles2,
  Alert,
} from '@grafana/ui';
import { User } from '../types';

import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';

const pageLimit = 20;

type UserListProps = {
  users: User[];
  userListError: any;
  onUserChecked: (event: React.FormEvent<HTMLInputElement>, userID: string) => void;
  checkedUsers: string[];
};

const UserList: React.FC<UserListProps> = ({ users, userListError, onUserChecked, checkedUsers }) => {
  const styles = useStyles2(getStyles);
  const [data, setData] = useState<User[] | undefined>(users);
  const [searchQuery, setSearchQuery] = useState('');
  const [isNoEmailUsersHidden, setIsNoEmailUsersHidden] = useState(false);

  const [paginationStates, setPaginationStates] = useState({
    totalPages: 1,
    currentPage: 1,
  });

  const renderUser = (user: User) => {
    return (
      <tr key={user.id}>
        <td className="width-2">
          <div className={styles.checkboxWrapper}>
            <Checkbox disabled={user.e_mail === ''} onChange={(event) => onUserChecked(event, user.id)} label="" />
          </div>
        </td>
        <td className="width-5">
          <div style={{ padding: '0px 8px' }}>{user.name}</div>
        </td>
        <td className="width-5">
          <div style={{ padding: '0px 8px' }} aria-label={user.e_mail?.length > 0 ? undefined : 'Empty email cell'}>
            {user.e_mail}
          </div>
        </td>
      </tr>
    );
  };

  const onIsNoEmailUsersHiddenChange = useCallback(
    (e) => {
      setIsNoEmailUsersHidden(e.currentTarget.checked);
    },
    [setIsNoEmailUsersHidden]
  );

  useEffect(() => {
    const getPaginatedUsers = (users: User[] | undefined) => {
      const offset = (paginationStates.currentPage - 1) * pageLimit;
      return users?.slice(offset, offset + pageLimit);
    };

    const filteredUser = users
      ?.filter((user) => {
        if (isNoEmailUsersHidden) {
          return user.e_mail !== '';
        }
        return true;
      })
      .filter((user) => {
        const match = user.name.toString().toLowerCase().indexOf(searchQuery.toLowerCase()) > -1;
        return match;
      });

    const data = getPaginatedUsers(filteredUser);
    const totalPages = Math.ceil(filteredUser.length / pageLimit);

    setPaginationStates((paginationStates) => ({
      ...paginationStates,
      currentPage: paginationStates.currentPage > totalPages ? 1 : paginationStates.currentPage,
      totalPages: filteredUser ? totalPages : 1,
    }));

    setData(data);
  }, [users, isNoEmailUsersHidden, searchQuery, paginationStates.currentPage]);

  return (
    <>
      <div className="page-action-bar">
        <FieldSet label="Selected Users">
          {checkedUsers.length > 0 ? (
            <HorizontalGroup wrap={true} style={{ marginBottom: '25px' }} align="flex-start" justify="flex-start">
              {checkedUsers.map((userID) => {
                const user = users.find((user) => user.id === userID);
                return <Tag key={userID} icon="user" name={`${user?.name} <${user?.e_mail}>`} />;
              })}
            </HorizontalGroup>
          ) : (
            <EmptySearchResult>You have not selected any user yet</EmptySearchResult>
          )}
        </FieldSet>
      </div>

      <div className="page-action-bar">
        <div className="gf-form gf-form--grow">
          <FieldSet label="Select Group members">
            <HorizontalGroup spacing="md" width="auto">
              <InlineField grow={true}>
                <Input
                  prefix={<Icon name="search" />}
                  suffix={<Icon name="trash-alt" onClick={() => setSearchQuery('')} />}
                  id="search-query"
                  name="search-query"
                  placeholder="Search users"
                  onChange={(e: any) => setSearchQuery(e.target.value)}
                />
              </InlineField>
              <InlineField grow={true} label="Hide users without email" transparent={true}>
                <InlineSwitch
                  checked={isNoEmailUsersHidden}
                  transparent={true}
                  onChange={onIsNoEmailUsersHiddenChange}
                ></InlineSwitch>
              </InlineField>
            </HorizontalGroup>
          </FieldSet>
        </div>
      </div>

      <div className="admin-list-table">
        <VerticalGroup spacing="md">
          {userListError && (
            <Alert title="User list error" severity="error">
              You must select at least one user
            </Alert>
          )}
          <table className="filter-table filter-table--hover form-inline">
            <thead>
              <tr>
                <th style={{ width: '1%' }} />
                <th>Name</th>
                <th>Email</th>
              </tr>
            </thead>
            <tbody>{data?.map((user) => renderUser(user))}</tbody>
          </table>
          <HorizontalGroup justify="center">
            <Pagination
              onNavigate={(page) =>
                setPaginationStates((paginationStates) => ({
                  ...paginationStates,
                  currentPage: page,
                }))
              }
              currentPage={paginationStates.currentPage}
              numberOfPages={paginationStates.totalPages}
              hideWhenSinglePage={true}
            />
          </HorizontalGroup>
        </VerticalGroup>
      </div>
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  checkboxWrapper: css`
    label {
      line-height: 1.2;
    }
  `,
});

export default UserList;
