import React, { useEffect, useState } from 'react';
import { Button, Field, FieldSet, FilterInput, Form, Input, VerticalGroup } from '@grafana/ui';
import { Page } from 'components/common/Page';
import { ROUTES, NAVIGATION_TITLE, NAVIGATION_SUBTITLE } from '../constants';
import { prefixRoute } from 'utils';
import { useDatasourceID } from 'hooks/useDatasourceID';
import { useQuery } from 'react-query';
import { getUsers } from 'api/getUsers.api';
import { User } from 'types';

interface ReportGroupType {
  name: string;
  description: string;
}

const CreateReportGroup = () => {
  const datasourceID = useDatasourceID();

  const [searchQuery, setSearchQuery] = useState('');

  const { data: users, error } = useQuery<User[], Error>(`users-${datasourceID}`, () => getUsers(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    retry: 0,
  });

  useEffect(() => {
    console.log('error', error);
  }, [error]);

  const renderUser = (user: User) => {
    return (
      <tr key={user.id}>
        <td className="width-4"></td>
        <td className="width-4">
          <div style={{ padding: '0px 8px' }}>{user.name}</div>
        </td>
        <td className="width-4">
          <div style={{ padding: '0px 8px' }} aria-label={user.e_mail?.length > 0 ? undefined : 'Empty email cell'}>
            {user.e_mail}
          </div>
        </td>
        <td className="width-4"></td>
      </tr>
    );
  };

  const submitCreateReportGroup = (formModel: ReportGroupType) => {};

  return (
    <Page
      headerProps={{
        title: NAVIGATION_TITLE,
        subTitle: NAVIGATION_SUBTITLE,
        backButton: {
          icon: 'arrow-left',
          href: prefixRoute(ROUTES.REPORT_GROUP),
        },
      }}
    >
      <Page.Contents>
        <Form onSubmit={submitCreateReportGroup}>
          {({ register, errors }) => (
            <>
              <FieldSet label="New Report Group">
                <Field
                  invalid={!!errors.name}
                  error={errors.name && errors.name.message}
                  label="Name"
                  description="Name of the report group"
                >
                  <Input
                    {...register('name', { required: 'Report group name is required' })}
                    id="report-group-name"
                    width={60}
                  />
                </Field>

                <Field label="description" description="Description of the report group">
                  <Input {...register('description')} id="report-group-description" width={60} />
                </Field>
              </FieldSet>

              <div className="page-action-bar">
                <div className="gf-form gf-form--grow">
                  <FilterInput placeholder="Search users" value={searchQuery} onChange={setSearchQuery} />
                </div>
              </div>

              <div className="admin-list-table">
                <VerticalGroup spacing="md">
                  <table className="filter-table filter-table--hover form-inline">
                    <thead>
                      <tr>
                        <th />
                        <th>Name</th>
                        <th>Email</th>
                        <th style={{ width: '1%' }} />
                      </tr>
                    </thead>
                    <tbody>{users?.map((user) => renderUser(user))}</tbody>
                  </table>
                </VerticalGroup>
              </div>

              <div className="gf-form-button-row">
                <Button type="submit" variant="primary">
                  Create Report Group
                </Button>
              </div>
            </>
          )}
        </Form>
      </Page.Contents>
    </Page>
  );
};
export { CreateReportGroup };
