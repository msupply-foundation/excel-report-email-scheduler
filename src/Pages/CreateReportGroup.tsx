import React from 'react';
import { Button, Field, FieldSet, Form, Input, LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { Page } from 'components/common/Page';
import { ROUTES, NAVIGATION_TITLE, NAVIGATION_SUBTITLE } from '../constants';
import { prefixRoute } from 'utils';
import { useDatasourceID } from 'hooks/useDatasourceID';
import { useMutation, useQuery } from 'react-query';
import { getUsers } from 'api/getUsers.api';
import { ReportGroupType, User } from 'types';
import UserList from 'components/UserList';
import { Controller } from 'react-hook-form';
import { createReportGroup, getReportGroupByID, getReportGroupMembersByGroupID } from 'api/ReportGroup';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';

const CreateReportGroup = ({ history, match }: any) => {
  const style = useStyles2(getStyles);
  const datasourceID = useDatasourceID();
  const { id: reportGroupIdToEdit } = match.params;
  const isEditMode = !!reportGroupIdToEdit;
  const [ready, setReady] = React.useState(false);

  const { data: users } = useQuery<User[], Error>(`users-${datasourceID}`, () => getUsers(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    refetchOnWindowFocus: false,
    retry: 0,
  });

  const { data: defaultReportGroup, isLoading: isReportGroupByIDLoading } = useQuery<ReportGroupType, Error>(
    `report-group-${reportGroupIdToEdit}`,
    () => getReportGroupByID(reportGroupIdToEdit),
    {
      enabled: isEditMode,
    }
  );

  const { data: defaultReportGroupMembers, isLoading: isReportGroupMembersByGroupIDLoading } = useQuery<any, Error>(
    `report-group-members-${reportGroupIdToEdit}`,
    () => getReportGroupMembersByGroupID(reportGroupIdToEdit),
    {
      enabled: isEditMode,
    }
  );

  React.useEffect(() => {
    if (
      !isReportGroupByIDLoading &&
      !isReportGroupMembersByGroupIDLoading &&
      !!defaultReportGroup &&
      !!defaultReportGroupMembers
    ) {
      defaultReportGroup.members = defaultReportGroupMembers.map((member: any) => member.userID);
      setReady(true);
    }
  }, [defaultReportGroup, defaultReportGroupMembers, isReportGroupByIDLoading, isReportGroupMembersByGroupIDLoading]);

  const createReportGroupMutation = useMutation((newReportGroup: ReportGroupType) => createReportGroup(newReportGroup));

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
          defaultValues={!!defaultReportGroup ? defaultReportGroup : undefined}
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
                    render={({ field: { onChange, value: selectedMembers } }) => {
                      console.log('selectedMembers', defaultReportGroup);
                      return (
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
                      );
                    }}
                    name="members"
                    control={control}
                    defaultValue={!!defaultReportGroup ? defaultReportGroup.members : []}
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

const getStyles = (theme: GrafanaTheme2) => ({
  loadingWrapper: css`
    display: flex;
    height: 50vh;
    align-items: center;
    justify-content: center;
  `,
});

export { CreateReportGroup };
