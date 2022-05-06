import React from 'react';
import { Controller } from 'react-hook-form';
import { Alert, Button, Field, FieldSet, Form, Input, PageToolbar, ToolbarButton } from '@grafana/ui';
import { useMutation, useQuery } from 'react-query';

import { Loading, Page, UserList } from '../../components';
import { ROUTES, NAVIGATION_TITLE, NAVIGATION_SUBTITLE, PLUGIN_BASE_URL } from '../../constants';
import { prefixRoute } from '../../utils';
import { useDatasourceID } from '../../hooks';
import { getUsers, createReportGroup, getReportGroupByID } from '../../api';
import { ReportGroupType, ReportGroupTypeWithMembersDetail, User } from '../../types';

const defaultFormValues: ReportGroupType = {
  id: '',
  name: '',
  description: '',
  members: [],
};

const CreateReportGroup = ({ history, match }: any) => {
  const datasourceID = useDatasourceID();
  const { id: reportGroupIdToEdit } = match.params;
  const isEditMode = !!reportGroupIdToEdit;
  const [ready, setReady] = React.useState(false);

  const { data: users, isLoading: isUsersLoading } = useQuery<User[], Error>(
    `users-${datasourceID}`,
    () => getUsers(datasourceID),
    {
      enabled: !!datasourceID,
      refetchOnMount: true,
      refetchOnWindowFocus: false,
      retry: 0,
    }
  );

  const [defaultReportGroup, setDefaultReportGroup] = React.useState<ReportGroupType>(defaultFormValues);

  const {
    data: defaultReportGroupFetched,
    isLoading: isReportGroupByIDLoading,
    isRefetching,
  } = useQuery<ReportGroupTypeWithMembersDetail, Error>(
    `report-group-${reportGroupIdToEdit}`,
    () => getReportGroupByID(reportGroupIdToEdit),
    {
      enabled: isEditMode,
      refetchOnMount: true,
      refetchOnWindowFocus: false,
      onError: () => {
        history.push(`${PLUGIN_BASE_URL}/report-groups/`);
        return;
      },
    }
  );

  React.useEffect(() => {
    if (!!defaultReportGroupFetched) {
      setDefaultReportGroup({
        ...defaultReportGroupFetched,
        members: defaultReportGroupFetched.members.map((member: User) => member.id),
      });
    }
  }, [defaultReportGroupFetched]);

  React.useEffect(() => {
    if (!isEditMode && !isUsersLoading) {
      setReady(true);
    }

    if (isEditMode && !isUsersLoading && !isRefetching && !isReportGroupByIDLoading) {
      setReady(true);
    }
  }, [isEditMode, isRefetching, isReportGroupByIDLoading, isUsersLoading]);

  const createReportGroupMutation = useMutation(
    (newReportGroup: ReportGroupType) => createReportGroup(newReportGroup),
    {
      onSuccess: () => {
        history.push(`${PLUGIN_BASE_URL}/report-groups/`);
        return;
      },
    }
  );

  const submitCreateReportGroup = (data: ReportGroupType): any => {
    createReportGroupMutation.mutate(data);
  };

  if (!ready) {
    return <Loading text="Loading..." />;
  }

  return (
    <Page
      headerProps={{
        title: NAVIGATION_TITLE,
        subTitle: NAVIGATION_SUBTITLE,
      }}
    >
      <Page.Contents>
        <PageToolbar
          parent="Report Groups"
          titleHref="#"
          parentHref={prefixRoute(ROUTES.REPORT_GROUP)}
          title={`${isEditMode ? 'Edit "' + defaultReportGroup?.name + '"' : 'New'}`}
          onGoBack={() => history.push(prefixRoute(ROUTES.REPORT_GROUP))}
        />
        <Form
          onSubmit={submitCreateReportGroup}
          validateOnMount={false}
          validateOn="onSubmit"
          defaultValues={defaultReportGroup}
          style={{
            marginTop: '30px',
          }}
        >
          {({ register, errors, control }) => {
            return (
              <>
                <FieldSet label="Details">
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

                {users ? (
                  <Controller
                    render={({ field: { onChange, value: selectedMembers } }) => (
                      <UserList
                        users={users}
                        userListError={errors.members}
                        checkedUsers={selectedMembers}
                        onUserChecked={(event, userID) => {
                          const updatedSelectedMembers = selectedMembers.includes(userID)
                            ? selectedMembers.filter((el) => el !== userID)
                            : [...selectedMembers, userID];

                          onChange(updatedSelectedMembers);
                        }}
                      />
                    )}
                    name="members"
                    control={control}
                  />
                ) : (
                  <Alert title="User(s) not found" severity="warning">
                    Report group must have members to be assigned from mSupply user. Please make sure you have mSupply
                    datasource selected in Plugin configuration.
                  </Alert>
                )}

                <div className="gf-form-button-row">
                  <Button type="submit" variant="primary">
                    {isEditMode ? 'Update' : 'Create'} Report Group
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
