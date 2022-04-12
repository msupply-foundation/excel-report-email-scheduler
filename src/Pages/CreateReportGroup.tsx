import React, { useEffect, useState } from 'react';
import { Button, Field, FieldSet, Input } from '@grafana/ui';
import { Page } from 'components/common/Page';
import { ROUTES, NAVIGATION_TITLE, NAVIGATION_SUBTITLE } from '../constants';
import { prefixRoute } from 'utils';
import { useDatasourceID } from 'hooks/useDatasourceID';
import { useQuery } from 'react-query';
import { getUsers } from 'api/getUsers.api';
import { User } from 'types';
import UserList from 'components/UserList';
import { useForm } from 'react-hook-form';
import { getBackendSrv } from '@grafana/runtime';

interface ReportGroupType {
  name: string;
  description?: string;
  selectedUsers: string[];
}

const CreateReportGroup = () => {
  const datasourceID = useDatasourceID();

  const { data: users, error } = useQuery<User[], Error>(`users-${datasourceID}`, () => getUsers(datasourceID), {
    enabled: !!datasourceID,
    refetchOnMount: true,
    retry: 0,
  });

  const {
    register,
    formState: { errors },
    handleSubmit,
    setValue,
  } = useForm({
    shouldFocusError: false,
  });

  const [checkedUsers, setCheckedUsers] = useState<string[]>([]);

  useEffect(() => {
    console.log('error', error);
  }, [error]);

  const onUserChecked = (event: React.FormEvent<HTMLInputElement>, userID: string) => {
    setCheckedUsers((oldArray: string[]) => {
      if (oldArray.includes(userID)) {
        return oldArray.filter((el) => el !== userID);
      } else {
        return [...oldArray, userID];
      }
    });
  };

  const submitCreateReportGroup = (data: ReportGroupType): any => {
    try {
      return getBackendSrv().post(
        `./api/plugins/msupplyfoundation-excelreportemailscheduler-datasource/resources/report-group`,
        data
      );
    } catch (error) {
      console.log(error);
    }
  };

  useEffect(() => {
    setValue('selectedUsers', checkedUsers);
  }, [checkedUsers, setValue]);

  useEffect(() => {
    register('selectedUsers', { validate: (value) => value.length > 0 });
  }, [register]);

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
        <form onSubmit={handleSubmit(submitCreateReportGroup)}>
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
            <UserList
              users={users}
              userListError={errors.selectedUsers}
              checkedUsers={checkedUsers}
              onUserChecked={onUserChecked}
            ></UserList>
          )}

          <div className="gf-form-button-row">
            <Button type="submit" variant="primary">
              Create Report Group
            </Button>
          </div>
        </form>
      </Page.Contents>
    </Page>
  );
};
export { CreateReportGroup };
