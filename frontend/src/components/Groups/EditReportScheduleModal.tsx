import React, { FC, useState, useEffect } from 'react';
import intl from 'react-intl-universal';
import { Modal, Button, ConfirmModal, Spinner } from '@grafana/ui';

import { queryCache, useMutation } from 'react-query';
import { deleteSchedule, sendTestEmail, updateSchedule } from 'api';

import { css } from 'emotion';

import { PanelList } from '../Schedules/PanelList';
import { useToggle } from 'hooks';
import { Schedule } from 'common/types';
import { EditScheduleForm } from '../Schedules/EditScheduleForm';
import { ScheduleKey } from 'common/enums';

type Props = {
  onClose: () => void;
  isOpen: boolean;
  reportSchedule: Schedule;
};

const modalAdjustments = css`
  top: 0;
  bottom: 0;
  width: 80%;
`;

const headerAdjustments = css`
  display: flex;
  flex: 1;
  justify-content: flex-end;
`;

export const EditReportScheduleModal: FC<Props> = ({ reportSchedule, onClose, isOpen }) => {
  const [schedule, setReportSchedule] = useState<Schedule>(reportSchedule);
  const [deleteAlertIsOpen, setDeleteAlertIsOpen] = useToggle(false);
  const [testEmails, { isLoading }] = useMutation(sendTestEmail);
  const [updateReportSchedule] = useMutation(updateSchedule, {
    onSuccess: () => queryCache.refetchQueries(['reportSchedules']),
  });

  const [deleteReportSchedule] = useMutation(deleteSchedule, {
    onSuccess: () => queryCache.refetchQueries(['reportSchedules']),
  });

  useEffect(() => {
    if (!schedule) {
      setReportSchedule(reportSchedule);
    }
  }, [schedule, reportSchedule]);

  const onUpdateSchedule = (key: ScheduleKey, newValue: string | number) => {
    const newState: Schedule = { ...schedule, [key]: newValue };
    setReportSchedule(newState);
    updateReportSchedule(newState);
  };

  const onConfirmDelete = () => {
    deleteReportSchedule(schedule);
    setDeleteAlertIsOpen();
    onClose();
  };

  return (
    <Modal
      className={modalAdjustments}
      onClickBackdrop={() => {}}
      title={intl.get('edit_report_schedule')}
      isOpen={isOpen}
      onDismiss={onClose}
    >
      <div className={headerAdjustments}>
        <EditScheduleForm schedule={schedule} onUpdate={onUpdateSchedule} />
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <Button size="md" variant="destructive" onClick={setDeleteAlertIsOpen}>
            {intl.get('delete')}
          </Button>
          {isLoading ? (
            <Spinner />
          ) : (
            <Button size="md" variant="primary" style={{ marginTop: '10px' }} onClick={() => testEmails(schedule.id)}>
              {intl.get('send_test_emails')}
            </Button>
          )}
        </div>
      </div>

      <PanelList schedule={reportSchedule} />

      <ConfirmModal
        isOpen={deleteAlertIsOpen}
        title={intl.get('delete_report_schedule')}
        body={intl.get('delete_report_schedule_question')}
        confirmText={intl.get('delete')}
        icon="exclamation-triangle"
        onConfirm={onConfirmDelete}
        onDismiss={setDeleteAlertIsOpen}
      />
    </Modal>
  );
};
