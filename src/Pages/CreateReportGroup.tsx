import React, { useEffect } from 'react';
import { Button, Field, FieldSet, Form, Input } from '@grafana/ui';
import { Page } from 'components/common/Page';
import { ROUTES, NAVIGATION_TITLE, NAVIGATION_SUBTITLE } from '../constants';
import { prefixRoute } from 'utils';
import { useDatasourceID } from 'hooks/useDatasourceID';
import { useMutation, useQuery } from 'react-query';
import { getUsers } from 'api/getUsers.api';
import { ReportGroupType, User } from 'types';
import UserList from 'components/UserList';
import { Controller } from 'react-hook-form';
import { createReportGroup } from 'api/ReportGroup';

const CreateReportGroup = () => {
  const datasourceID = useDatasourceID();

  const { data: users, error } = useQuery<User[], Error>(`users-${datasourceID}`, () => getUsers(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const createReportGroupMutation = useMutation(
    (newReportGroup: ReportGroupType) => createReportGroup(newReportGroup),
    {}
  );

  useEffect(() => {
    console.log('error', error);
  }, [error]);

  const submitCreateReportGroup = (data: ReportGroupType): any => {
    console.log('data', data);
    createReportGroupMutation.mutate(data);
  };

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
        <Form onSubmit={submitCreateReportGroup} validateOnMount={false} validateOn="onSubmit">
          {({ register, errors, control }) => {
            return (
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

                {users && (
                  <Controller
                    render={({ field: { onChange, value: selectedMembers } }) => {
                      return (
                        <UserList
                          users={users}
                          userListError={errors.selectedUsers}
                          checkedUsers={selectedMembers}
                          onUserChecked={(event, userID) => {
                            const updatedSelectedMembers = selectedMembers.includes(userID)
                              ? selectedMembers.filter((el) => el !== userID)
                              : [...selectedMembers, userID];

                            onChange(updatedSelectedMembers);
                          }}
                        ></UserList>
                      );
                    }}
                    name="selectedUsers"
                    control={control}
                    defaultValue={[]}
                  />
                )}

                <div className="gf-form-button-row">
                  <Button type="submit" variant="primary">
                    Create Report Group
                  </Button>
                </div>
              </>
            );
          }}
        </Form>
      </Page.Contents>
    </Page>
  );
};
export { CreateReportGroup };
