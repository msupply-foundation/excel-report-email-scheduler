import React from 'react';
import { Button, Field, FieldSet, Form, Input, LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { Page } from 'components/common';
import { ROUTES, NAVIGATION_TITLE, NAVIGATION_SUBTITLE, PLUGIN_BASE_URL } from '../../constants';
import { prefixRoute } from 'utils';
import { useDatasourceID } from 'hooks/useDatasourceID';
import { useMutation, useQuery } from 'react-query';
import { getUsers } from 'api/getUsers.api';
import { ReportGroupType, ReportGroupTypeWithMembersDetail, User } from 'types';
import UserList from 'components/UserList';
import { Controller } from 'react-hook-form';
import { createReportGroup, getReportGroupByID } from 'api/ReportGroup';
import { css } from '@emotion/css';

const defaultFormValues: ReportGroupType = {
  id: '',
  name: '',
  description: '',
  members: [],
};

const CreateReportGroup = ({ history, match }: any) => {
  const style = useStyles2(getStyles);
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
        console.log('On success called');
        history.push(`${PLUGIN_BASE_URL}/report-groups/`);
        return;
      },
    }
  );

  const submitCreateReportGroup = (data: ReportGroupType): any => {
    createReportGroupMutation.mutate(data);
  };

  if (!ready) {
    return (
      <div className={style.loadingWrapper}>
        <LoadingPlaceholder text="Loading..." />
      </div>
    );
  }

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
        <Form
          onSubmit={submitCreateReportGroup}
          validateOnMount={false}
          validateOn="onSubmit"
          defaultValues={defaultReportGroup}
        >
          {({ register, errors, control }) => {
            return (
              <>
                <FieldSet label={`${isEditMode ? 'Edit "' + defaultReportGroup?.name + '"' : 'New'} Report Group`}>
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
                      ></UserList>
                    )}
                    name="members"
                    control={control}
                  />
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

const getStyles = () => ({
  loadingWrapper: css`
    display: flex;
    height: 50vh;
    align-items: center;
    justify-content: center;
  `,
});

export { CreateReportGroup };
