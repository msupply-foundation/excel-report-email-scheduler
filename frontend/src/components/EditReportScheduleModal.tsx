import React, { FC, useState, useEffect } from 'react';
import intl from 'react-intl-universal';
import { Input, Field, Modal, Button, Select, ConfirmModal } from '@grafana/ui';

import { queryCache, useMutation, useQuery } from 'react-query';
import {
  deleteSchedule,
  deleteReportContent,
  createReportContent,
  getPanels,
  getReportContent,
  updateSchedule,
} from 'api';

import { css } from 'emotion';

import { PanelList } from './PanelList';
import { useToggle } from 'hooks';
import { ReportContent, Schedule } from 'common/types';
import { EditScheduleForm } from './Schedules/EditScheduleForm';
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

  const { data: content } = useQuery<ReportContent[] | null[]>({
    queryKey: ['reportContent', reportSchedule?.id],
    queryFn: getReportContent,
    config: { enabled: !!reportSchedule },
  });

  const { data: panels } = useQuery({
    queryKey: ['panels'],
    queryFn: getPanels,
    config: { enabled: !!reportSchedule },
  });

  const [updateReportSchedule] = useMutation(updateSchedule, {
    onSuccess: () => queryCache.refetchQueries(['reportSchedules']),
  });

  const [createContent] = useMutation(createReportContent, {
    onSuccess: () => queryCache.refetchQueries(['reportContent', reportSchedule?.id]),
  });

  const [deleteContent] = useMutation(deleteReportContent, {
    onSuccess: () => queryCache.refetchQueries(['reportContent', reportSchedule?.id]),
  });

  const [deleteReportSchedule] = useMutation(deleteSchedule, {
    onSuccess: () => queryCache.refetchQueries(['reportSchedules']),
  });

  const onTogglePanel = async (panel: any) => {
    // Exists, need to remove.
    if (!!content?.[panel?.id]) {
      await deleteContent(panel);
    } else {
      if (content) {
        await createContent({ scheduleID: schedule?.id, panelID: panel?.id });
      }
    }
  };

  useEffect(() => {
    if (!schedule) setReportSchedule(reportSchedule);
  }, [reportSchedule]);

  // TODO: Handle error cases
  const onUpdateSchedule = (key: ScheduleKey, newValue: string | number) => {
    const newState: Schedule = { ...schedule, [key]: newValue };
    // Optimistically update state to reflect changes immediately in UI.
    setReportSchedule(newState);
    updateReportSchedule(newState);
  };

  const onConfirmDelete = () => {
    deleteReportSchedule(schedule);
    setDeleteAlertIsOpen();
    onClose();
  };

  const onUpdateReportContent = (key: string, newValue: any) => {
    if (key === 'storeID') {
      const csv = newValue
        .map((store: any) => {
          return store.id;
        })
        .join(', ');
      console.log(csv);
    }
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
        <Button size="md" variant="destructive" onClick={setDeleteAlertIsOpen}>
          {intl.get('delete')}
        </Button>
      </div>

      <PanelList schedule={reportSchedule} />

      <ConfirmModal
        isOpen={deleteAlertIsOpen}
        title={intl.get('delete_report_group')}
        body={intl.get('delete_report_group_question')}
        confirmText={intl.get('delete')}
        icon="exclamation-triangle"
        onConfirm={onConfirmDelete}
        onDismiss={setDeleteAlertIsOpen}
      />
    </Modal>
  );
};
