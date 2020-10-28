import React, { FC, useState, useEffect } from 'react';
import intl from 'react-intl-universal';
import { Input, Field, Modal, Button, Select, ConfirmModal } from '@grafana/ui';

import { queryCache, useMutation, useQuery } from 'react-query';
import { deleteReportContent, createReportContent, getPanels, getReportContent, updateSchedule } from 'api';

import { css } from 'emotion';

import { PanelList } from './PanelList';
import { useToggle } from 'hooks';
import { Schedule } from 'common/types';
import { EditScheduleForm } from './Schedules/EditScheduleForm';

type Props = {
  onClose: () => void;
  isOpen: boolean;
  reportSchedule: ReportSchedule;
};

type ReportContent = {
  id?: string;
  panelID?: string;
};

type ReportSchedule = {
  id?: string;
  interval?: number;
  nextReportTime?: number;
  name?: string;
  description?: string;
  lookback?: number;
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
  const [schedule, setReportSchedule] = useState<ReportSchedule>(reportSchedule);
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

  const [createContent] = useMutation(createReportContent, {
    onSuccess: () => queryCache.refetchQueries(['reportContent', reportSchedule?.id]),
  });
  const [deleteContent] = useMutation(deleteReportContent, {
    onSuccess: () => queryCache.refetchQueries(['reportContent', reportSchedule?.id]),
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

  const onUpdateSchedule = (key: any, newValue: string) => {
    const newState: Schedule = { ...schedule, [key]: newValue };

    // Optimistically update state
    // TODO: Handle error case
    setReportSchedule(newState);
    updateSchedule(newState);
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
        <EditScheduleForm schedule={schedule} onUpdate={() => {}} />
        <Button size="md" variant="destructive">
          {intl.get('delete')}
        </Button>
      </div>

      <PanelList
        scheduleID={reportSchedule?.id}
        checked={content}
        withChecks
        onRowPress={onTogglePanel}
        data={panels}
        titleKey="title"
        descriptionKey="description"
      />

      <ConfirmModal
        isOpen={deleteAlertIsOpen}
        title={intl.get('delete_report_group')}
        body={intl.get('delete_report_group_question')}
        confirmText={intl.get('delete')}
        icon="exclamation-triangle"
        onConfirm={() => {}}
        onDismiss={setDeleteAlertIsOpen}
      />
    </Modal>
  );
};
